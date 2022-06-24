package sql

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/edgelesssys/ego/ecrypto"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/go-sql-driver/mysql"
)

/*
   The Obscuro Enclave (OE) needs a way to persist data into a trusted database. Trusted not to reveal that data to anyone but that particular enclave.

   To achieve this, the OE must first perform Remote Attestation (RA), which gives it confidence that it is connected to
	a trusted version of software running on trusted hardware. The result of this process is a Certificate which can be
	used to set up a trusted TLS connection into the database.

   The next step is to configure the database schema and users in such a way that the OE knows that the db engine will
	only allow itself access to it. This is achieved by creating a "Manifest" file that contains the SQL init code and a
	DBClient Certificate that is known only to the OE.

	This "DBClient" Cert is used by the database to authenticate that it is communicating to the entity that has initialised that schema.

	--------

	In more detail :
   The edgeless DB is unusable when it's first started, as an owner, we must initially:
     - do a remote attestation on the report provided by {edbAddress}/quote
     - create a ca cert that will authenticate our DB users going forwards
     - prepare a manifest.json that contains that CA cert and some SQL to initialise the DB tables and user
     - submit that manifest.json file to {edbAddress}/manifest, using the certificate provided from /quote to authenticate
     - seal and persist the manifest.json and the certs so we can reconnect if enclave is restarted

   When connecting to an already-initialized edgeless DB we must:
     - perform remote attestation on the edgeless db
     - unseal the manifest.json and get the hash of it, also unseal the certificate that edb was initialized with
	 - verify the /signature request for edgeless DB matches the manifest.json hash
     - connect to edb with the persisted cert - it's now safe to read and write to the DB

	Some useful documentation for this:
		Main edb docs: https://docs.edgeless.systems/edgelessdb/#/
		EDB demo docs: https://github.com/edgelesssys/edgelessdb/tree/main/demo
		// Note: due to an issue with the dependency, I've duplicated the relevant parts of the ERA tool into edb_attestation.go
		ERA - remote attestation tool: https://github.com/edgelesssys/era
*/

const (
	edbHTTPPort          = "8080"
	edbManifestEndpoint  = "/manifest"
	edbSignatureEndpoint = "/signature"

	dataDir         = "/data"
	certIssuer      = "obscuroCA"
	certSubject     = "obscuroUser"
	enclaveHostName = "enclave"

	dbUser    = "obscuro"
	dbName    = "obsdb"
	tableName = "keyvalue"
	keyCol    = "ky"
	valueCol  = "val"

	// The attestation config comes from here (https://github.com/edgelesssys/edgelessdb/releases/latest/download/edgelessdb-sgx.json)
	//     todo: revisit whether we want this hardcoded
	edbAttestationConf = "{\n\t\"SecurityVersion\": 2,\n\t\"ProductID\": 16,\n\t\"SignerID\": \"67d7b00741440d29922a15a9ead427b6faf1d610238ae9826da345cea4fee0fe\"\n}"

	// change this flag to true to debug issues with edgeless DB (and start EDB process with -e EDG_EDB_DEBUG=1
	//   this will give you:
	//  	- verbose logging on EDB
	//		- write the edb.pem file out for connecting to Edgeless DB services manually
	//		- versions of files created with a '.unsealed' suffix that can be used to connect to the database using mysql-client
	debugMode = false
)

var (
	manifestFilepath       = filepath.Join(dataDir, "manifest.json")
	caCertFilepath         = filepath.Join(dataDir, "ca-cert.pem")
	userCertFilepath       = filepath.Join(dataDir, "user-cert.pem")
	userKeyFilepath        = filepath.Join(dataDir, "user-key.pem")
	attestationCfgFilepath = filepath.Join(dataDir, "edgelessdb-sgx.json")

	manifestSQLStatements = []string{
		fmt.Sprintf("CREATE USER %s REQUIRE ISSUER '/CN=%s' SUBJECT '/CN=%s'", dbUser, certIssuer, certSubject),
		fmt.Sprintf("CREATE DATABASE %s", dbName),
		fmt.Sprintf("CREATE TABLE %s.%s (%s varbinary(64) primary key, %s blob)", dbName, tableName, keyCol, valueCol),
		fmt.Sprintf("GRANT ALL ON %s.%s TO %s", dbName, tableName, dbUser),
	}
)

type manifest struct {
	SQL   []string `json:"sql"`
	Cert  string   `json:"ca"`
	Debug bool     `json:"debug"`
}

type EdgelessDBConfig struct {
	Host string
}

func EdgelessDBConnector(edbCfg EdgelessDBConfig) (ethdb.Database, error) {
	// Before we try to connect to the Edgeless DB we have to do Remote Attestation (RA) on it
	// the RA will ensure that we are connecting to a database that will not leak any data.
	// The RA will return a Certificate which we'll use for the TLS mutual authentication when we connect to the database.
	// The trust path is as follows:
	// 1. The Obscuro Enclave performs RA on the database enclave, and the RA object contains a certificate which only the database enclave controls.
	// 2. Connecting to the database via mutually authenticated TLS using the above certificate, will give the Obscuro enclave confidence that it is only giving data away to some code and hardware it trusts.
	edbPEM, err := performEDBRemoteAttestation(edbCfg.Host)
	if err != nil {
		return nil, fmt.Errorf("remote attestation of edgeless DB failed - %w", err)
	}

	// client used to make secure HTTP requests to Edgeless DB using the ca-cert we have received
	edbHTTPClient, err := prepareEDBHTTPClient(edbPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare http client from EdgelessDB cert PEM - %w", err)
	}

	// IF DEBUGGING: it can be useful to persist the edb.pem so you can connect to edb from mysql-client on the container
	if debugMode {
		err = sealAndPersist(edbPEM, filepath.Join(dataDir, "edb.pem"))
		if err != nil {
			return nil, fmt.Errorf("failed to writeout edb.pem file for debugging - %w", err)
		}
	}

	// now we know we are talking to a secure enclave, we can get the manifest and connect (or initialise if first time)
	manifest, found, err := readManifestIfExists()
	if err != nil {
		// this doesn't happen if the manifest file just didn't exist, maybe there was an IO error
		return nil, fmt.Errorf("failed to read manifest file - %w", err)
	}
	if !found {
		// this is the first time we have connected to this EDB, we will create certificates and a manifest to initialise it
		log.Info("No manifest found, creating one and initializing edb")
		manifest, err = createManifestAndInitEDB(edbCfg.Host, edbHTTPClient)
		if err != nil {
			return nil, err
		}

		// Note: it usually takes around 10-15 seconds for edb to initialise and restart
		log.Info("Waiting 30 seconds for EDB restart after initialization...")
		time.Sleep(30 * time.Second)
	}

	// we check that this edgeless DB was initialized with the manifest we expected (which is only known to this enclave)
	log.Info("Validating edb signature against expected manifest...")
	err = verifyEdgelessDB(edbCfg.Host, manifest, edbHTTPClient)
	if err != nil {
		return nil, err
	}

	// connect to EDB (standard mysql-type connection, using certificate derived from the CA cert in the manifest)
	log.Info("Setting up SQL connection to edb...")
	edbSQL, err := connectToEdgelessDB(edbCfg.Host, edbPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to EdgelessDB - %w", err)
	}

	// wrap it in our eth-compatible key-value store layer
	return CreateSQLEthDatabase(edbSQL)
}

func prepareEDBHTTPClient(edbPEM string) (*http.Client, error) {
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM([]byte(edbPEM)); !ok {
		return nil, fmt.Errorf("failed to append to CA cert from edb cert pem")
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    caCertPool,
				MinVersion: tls.VersionTLS12,
			},
		},
	}, nil
}

// perform the SGX enclave attestation to verify edb running in a legit enclave and with expected edb version etc.
func performEDBRemoteAttestation(edbHost string) (string, error) {
	// we need to make sure this dir exists before we start read/writing files in there
	err := os.MkdirAll(dataDir, 0o644)
	if err != nil {
		return "", err
	}
	err = prepareEDBAttestationRequirementsConf(attestationCfgFilepath)
	if err != nil {
		return "", fmt.Errorf("failed to prepare latest edb attestation config file - %w", err)
	}
	log.Info("Verifying attestation from edgeless DB...")
	edbHTTPAddr := fmt.Sprintf("%s:%s", edbHost, edbHTTPPort)
	certs, tcbStatus, err := GetCertificate(edbHTTPAddr, attestationCfgFilepath)
	if err != nil {
		// todo should we check the error type with: err == attestation.ErrTCBLevelInvalid?
		// for now it's maximum strictness (we can revisit this and permit some tcbStatuses if desired)
		return "", fmt.Errorf("attestation failed, host=%s, tcbStatus=%s, err=%w", edbHTTPAddr, tcbStatus, err)
	}
	if len(certs) == 0 {
		return "", fmt.Errorf("no certificates found from edgeless db attestation process")
	}

	log.Info("Successfully verified edb attestation and retrieved certificate.")
	// the last cert in the list is the CA
	return string(pem.EncodeToMemory(certs[len(certs)-1])), nil
}

func prepareEDBAttestationRequirementsConf(filepath string) error {
	// This json blob provides confidence in the version of edgeless DB we are talking to.
	// The latest json for comparison is available here:
	//     https://github.com/edgelesssys/edgelessdb/releases/latest/download/edgelessdb-sgx.json
	return os.WriteFile(filepath, []byte(edbAttestationConf), 0o444)
}

func createManifestAndInitEDB(edbHost string, httpClient *http.Client) (*manifest, error) {
	caCert := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		IsCA:                  true,
		BasicConstraintsValid: true,
		// this subject must match the subject authorised in the manifest.json
		Subject:   pkix.Name{CommonName: certIssuer},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
		DNSNames:  []string{enclaveHostName},
	}
	caPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key for CA cert to init Edgeless DB - %w", err)
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create CA cert - %w", err)
	}
	caPEMBuf := new(bytes.Buffer)
	err = pem.Encode(caPEMBuf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ca cert pem - %w", err)
	}

	caPEM := caPEMBuf.String()
	manifest := &manifest{
		SQL:   manifestSQLStatements,
		Cert:  caPEM,
		Debug: debugMode,
	}
	err = initialiseEdgelessDB(edbHost, manifest, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise edgeless DB with created manifest - %w \nmanifest: %v", err, manifest)
	}

	// store certificates for DB connection
	err = prepareCertificates(caCert, caPEM, caPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare and persist certs for edb connection - %w", err)
	}

	// persist the manifest for any future restarts of the enclave
	err = writeManifest(manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to persist manifest file - %w", err)
	}

	return manifest, nil
}

// initialiseEdgelessDB sends a manifest over http to the edgeless DB with its initial config
func initialiseEdgelessDB(edbHost string, m *manifest, httpClient *http.Client) error {
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest json - %w", err)
	}
	url := fmt.Sprintf("https://%s:%s%s", edbHost, edbHTTPPort, edbManifestEndpoint)
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("faild to create manifest initialization req - %w", err)
	}

	_, err = executeHTTPReq(httpClient, req)
	if err != nil {
		return fmt.Errorf("manifest initialization req failed - %w", err)
	}

	return nil
}

// verifyEdgelessDB requests the /signature from the edb, it should match the hash of the manifest we expected
func verifyEdgelessDB(edbHost string, m *manifest, httpClient *http.Client) error {
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest to json - %w", err)
	}
	h := sha256.Sum256(b)
	expectedHash := hex.EncodeToString(h[:])

	url := fmt.Sprintf("https://%s:%s%s", edbHost, edbHTTPPort, edbSignatureEndpoint)
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return fmt.Errorf("faild to create edb signature req - %w", err)
	}

	edbHash, err := executeHTTPReq(httpClient, req)
	if err != nil {
		return fmt.Errorf("failed to receive edbHash from /signature request - %w", err)
	}
	if expectedHash != string(edbHash) {
		return fmt.Errorf("hash from edb /signature request didn't match expected hash of manifest.json, expected=%s, found=%s", expectedHash, edbHash)
	}
	log.Info("EDB signature matched the expected hash from our manifest (%s)", expectedHash)

	return nil
}

// create Go standard database/sql connection to edb using a mysql driver
func connectToEdgelessDB(edbHost string, edbPEM string) (*sql.DB, error) {
	caCertPool := x509.NewCertPool()

	if ok := caCertPool.AppendCertsFromPEM([]byte(edbPEM)); !ok {
		return nil, fmt.Errorf("failed to append edb cert to mysql CA cert pool")
	}

	userCert, err := readAndUnseal(userCertFilepath)
	if err != nil {
		return nil, err
	}
	userKey, err := readAndUnseal(userKeyFilepath)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(userCert, userKey)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare keypair from cert and key - %w", err)
	}
	err = mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to prepare certs for mysql connection - %w", err)
	}
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = edbHost
	cfg.User = dbUser
	cfg.DBName = dbName
	cfg.TLSConfig = "custom"
	dsn := cfg.FormatDSN()
	log.Info("Configuring mysql connection: %s", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mysql connection to edb - %w", err)
	}
	return db, nil
}

// readManifestIfExists returns manifest if it exists, whether it was found as a boolean, error
func readManifestIfExists() (*manifest, bool, error) {
	var manifest manifest
	_, err := os.Stat(manifestFilepath)
	if err != nil {
		if os.IsNotExist(err) {
			// we don't consider the file being missing as an error scenario, it's just not initialized
			return nil, false, nil
		}
		// failed to open file
		return nil, false, fmt.Errorf("failed to open manifest file - %w", err)
	}
	jsonData, err := readAndUnseal(manifestFilepath)
	if err != nil {
		return nil, false, err
	}
	err = json.Unmarshal(jsonData, &manifest)
	if err != nil {
		// failed to unmarshal the json
		return nil, false, fmt.Errorf("failed to unmarshal manifest json - %w", err)
	}
	log.Info("Successfully loaded manifest from disk.")
	return &manifest, true, nil
}

func writeManifest(m *manifest) error {
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest json - %w", err)
	}

	return sealAndPersist(string(b), manifestFilepath)
}

// prepareCertificates persists the ca-cert we generated for the manifest and creates and persists a user cert + key from it
func prepareCertificates(caCert *x509.Certificate, caPEM string, caPrivKey *ecdsa.PrivateKey) error {
	err := sealAndPersist(caPEM, caCertFilepath)
	if err != nil {
		return err
	}

	userCert := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		// the issuer and subject have to match those submitted in manifest.json
		Issuer:    pkix.Name{CommonName: certIssuer},
		Subject:   pkix.Name{CommonName: certSubject},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
	}
	certPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key for user cert - %w", err)
	}
	userCertBytes, err := x509.CreateCertificate(rand.Reader, userCert, caCert, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return fmt.Errorf("failed to prepare user certificate - %w", err)
	}

	userCertPEM := new(bytes.Buffer)
	err = pem.Encode(userCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: userCertBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to PEM encode user certificate - %w", err)
	}
	err = sealAndPersist(userCertPEM.String(), userCertFilepath)
	if err != nil {
		return err
	}

	certKeyPEM := new(bytes.Buffer)
	privKeyOut, err := x509.MarshalPKCS8PrivateKey(certPrivKey)
	if err != nil {
		return fmt.Errorf("failed to marshal cert priv key - %w", err)
	}
	err = pem.Encode(certKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyOut,
	})
	if err != nil {
		return fmt.Errorf("failed to pem encode the user private key - %w", err)
	}
	return sealAndPersist(certKeyPEM.String(), userKeyFilepath)
}

func sealAndPersist(contents string, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s - %w", filepath, err)
	}
	defer f.Close()

	if debugMode {
		fUnseal, _ := os.Create(filepath + ".unsealed")
		_, err = fUnseal.WriteString(contents)
		if err != nil {
			return err
		}
		_ = fUnseal.Close()
	}

	// todo: do we prefer to seal with product key for upgradability or unique key to require fresh db with every code change
	enc, err := ecrypto.SealWithProductKey([]byte(contents), nil)
	if err != nil {
		return fmt.Errorf("failed to seal contents bytes with enclave key to persist in %s - %w", filepath, err)
	}
	_, err = f.Write(enc)
	if err != nil {
		return fmt.Errorf("failed to write manifest json file - %w", err)
	}
	return nil
}

func readAndUnseal(filepath string) ([]byte, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s - %w", filepath, err)
	}

	data, err := ecrypto.Unseal(b, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to unseal data from file %s - %w", filepath, err)
	}
	return data, nil
}

// executeHTTPReq executes an HTTP request, returns an error if the response code was outside of 200-299, returns response body as bytes if there was a response body
func executeHTTPReq(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed url=%s - %w", req.URL.String(), err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var msg []byte
		_, err := resp.Body.Read(msg)
		if err != nil {
			return nil, fmt.Errorf("req failed url=%s, statusCode=%d, failed to read status text", req.URL.String(), resp.StatusCode)
		}
		return nil, fmt.Errorf("req failed url=%s status: %d %s", req.URL.String(), resp.StatusCode, msg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// success status code but no body, ignoring error
		return []byte{}, nil //nolint:nilerr
	}
	return body, nil
}
