package sql

import (
	"bytes"
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
	"github.com/edgelesssys/ego/ecrypto"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/go-sql-driver/mysql"
)

/*
   Getting an edgeless DB connection is a bit of a dance in order to achieve mutual trust.

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
		// Note: due to an issue with the dependency, I've duplicated the relevant parts of the ERA tool into edbattestation.go
		ERA - remote attestation tool: https://github.com/edgelesssys/era
*/

var manifestSQLStatements = []string{
	"CREATE USER obscuro REQUIRE ISSUER '/CN=obscuroCA' SUBJECT '/CN=obscuroUser'",
	"CREATE DATABASE obsdb",
	"CREATE TABLE obsdb.keyvalue (ky varbinary(64) primary key, val blob)",
	"GRANT ALL ON obsdb.keyvalue TO obscuro",
}

const (
	dataDir                = "/data/"
	manifestFilepath       = dataDir + "manifest.json"
	caCertFilepath         = dataDir + "ca-cert.pem"
	userCertFilepath       = dataDir + "user-cert.pem"
	userKeyFilepath        = dataDir + "user-key.pem"
	attestationCfgFilepath = dataDir + "edgelessdb-sgx.json"
	edbHttpPort            = "8080"
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
	// before we try to connect to the Edgeless DB we have to do remote attestation on it
	edbPEM, err := performRemoteAttestation(edbCfg.Host)
	if err != nil {
		return nil, fmt.Errorf("remote attestation of edgeless DB failed - %w", err)
	}
	log.Info("retrieved edb PEM: %s", edbPEM)

	// IF DEBUGGING: it can be useful to persist the edb.pem so you can connect to edb from mysql-client on the container
	//_ = sealAndPersist(edbPEM, dataDir+"edb.pem")

	// now we know we are talking to a secure enclave, we can get the manifest and connect (or initialise if first time)
	manifest, err := readManifestIfExists()
	if err != nil {
		// this doesn't happen if the manifest file just didn't exist, maybe there was an IO error
		return nil, fmt.Errorf("failed to read manifest file - %w", err)
	}
	if manifest == nil {
		// this is the first time we have connected to this EDB, we will create certificates and a manifest to initialise it
		log.Info("No manifest found, creating one and initializing edb")
		manifest, err = createManifestAndInitEDB(edbCfg.Host, edbPEM)
		if err != nil {
			return nil, err
		}

		// Note: it usually takes around 10-15 seconds for edb to initialise and restart
		log.Info("Waiting 30 seconds for EDB restart after initialization...")
		time.Sleep(30 * time.Second)
	}

	// we check that this edgeless DB was initialized with the manifest we expected (which is only known to this enclave)
	log.Info("Validating edb signature against expected manifest...")
	err = verifyEdgelessDB(edbCfg.Host, manifest, edbPEM)
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

// perform the SGX enclave attestation to verify edb running in a legit enclave and with expected edb version etc.
func performRemoteAttestation(edbHost string) (string, error) {
	// we need to make sure this dir exists before we start read/writing files in there
	err := os.MkdirAll(dataDir, 0644)
	if err != nil {
		return "", err
	}
	err = prepareEDBAttestationConf(attestationCfgFilepath)
	if err != nil {
		return "", fmt.Errorf("failed to prepare latest edb attestation config file - %w", err)
	}
	log.Info("Verifying attestation from edgeless DB...")
	edbHttpAddr := fmt.Sprintf("%s:%s", edbHost, edbHttpPort)
	certs, tcbStatus, err := GetCertificate(edbHttpAddr, attestationCfgFilepath)
	if err != nil {
		// todo should we check the error type with: err == attestation.ErrTCBLevelInvalid?
		// for now it's maximum strictness (we can revisit this and permit some tcbStatuses if desired)
		return "", fmt.Errorf("attestation failed, host=%s, tcbStatus=%s, err=%w", edbHttpAddr, tcbStatus, err)
	}
	if len(certs) == 0 {
		return "", fmt.Errorf("no certificates found from edgeless db attestation process")
	}

	log.Info("Successfully verified edb attestation and retrieved certificate.")
	return string(pem.EncodeToMemory(certs[0])), nil
}

func prepareEDBAttestationConf(filepath string) error {
	// This json blob provides confidence in the version of edgeless DB we are talking to.
	// The latest json for comparison is available here:
	//     https://github.com/edgelesssys/edgelessdb/releases/latest/download/edgelessdb-sgx.json
	//
	// Todo: revisit how we want to enforce this going forwards, but for now I'm just hardcoding the latest json blob:
	b := []byte("{\n\t\"SecurityVersion\": 2,\n\t\"ProductID\": 16,\n\t\"SignerID\": \"67d7b00741440d29922a15a9ead427b6faf1d610238ae9826da345cea4fee0fe\"\n}")
	return os.WriteFile(filepath, b, 0644)
}

func createManifestAndInitEDB(edbHost string, edbPEM string) (*manifest, error) {
	caCert := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		IsCA:                  true,
		BasicConstraintsValid: true,
		Subject:               pkix.Name{CommonName: "obscuroCA"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		DNSNames:              []string{"enclave"},
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
		SQL:  manifestSQLStatements,
		Cert: caPEM,
		// IF DEBUGGING: this should be set to enable verbose logging
		//     Note: it also requires EDG_EDB_DEBUG=1 on the process, see https://docs.edgeless.systems/edgelessdb/#/reference/configuration
		// Debug: true,
	}
	err = initialiseEdgelessDB(edbHost, manifest, edbPEM)
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
func initialiseEdgelessDB(edbHost string, m *manifest, edbPEM string) error {
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest json - %w", err)
	}
	url := fmt.Sprintf("https://%s:%s/manifest", edbHost, edbHttpPort)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("faild to create manifest initialization req - %w", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM([]byte(edbPEM)); !ok {
		return fmt.Errorf("failed to append to CA cert from edb cert pem")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("manifest initialization req failed - %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var msg []byte
		_, err := resp.Body.Read(msg)
		if err != nil {
			return fmt.Errorf("manifest initialization req failed with status code: %d, failed to read status text", resp.StatusCode)
		}
		return fmt.Errorf("manifest initialization req failed with status: %d %s", resp.StatusCode, msg)
	}
	return nil
}

// verifyEdgelessDB requests the /signature from the edb, it should match the hash of the manifest we expected
func verifyEdgelessDB(edbHost string, m *manifest, edbPEM string) error {
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest to json - %w", err)
	}
	h := sha256.Sum256(b)
	expectedHash := hex.EncodeToString(h[:])

	url := fmt.Sprintf("https://%s:%s/signature", edbHost, edbHttpPort)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("faild to create edb signature req - %w", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(edbPEM))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	resp, err := client.Do(req)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var msg []byte
		_, err := resp.Body.Read(msg)
		if err != nil {
			return fmt.Errorf("edb /signature req failed with status code: %d, failed to read status text", resp.StatusCode)
		}
		return fmt.Errorf("edb /signature req failed with status: %d %s", resp.StatusCode, msg)
	}
	var edbHash []byte
	edbHash, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read edbHash from /signature response - %w", err)
	}
	if expectedHash != string(edbHash) {
		return fmt.Errorf("hash from edb /signature request didn't match expected hash of manifest.json, expected=%s, found=%s resp=%v", expectedHash, edbHash, resp)
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

	cert, err := tls.LoadX509KeyPair(userCertFilepath, userKeyFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate from block - %w", err)
	}
	err = mysql.RegisterTLSConfig("custom", &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            caCertPool,
		Certificates:       []tls.Certificate{cert},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to prepare certs for mysql connection - %w", err)
	}
	cfg := mysql.NewConfig()
	cfg.Net = "tcp"
	cfg.Addr = edbHost
	cfg.User = "obscuro"
	cfg.TLSConfig = "custom"
	cfg.DBName = "obsdb"
	dsn := cfg.FormatDSN()
	log.Info("Configuring mysql connection: %s", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mysql connection to edb - %w", err)
	}
	return db, nil
}

func readManifestIfExists() (*manifest, error) {
	var manifest manifest
	_, err := os.Stat(manifestFilepath)
	if err != nil {
		if os.IsNotExist(err) {
			// we don't consider the file being missing as an error scenario, it's just not initialized
			return nil, nil
		}
		// failed to open file
		return nil, fmt.Errorf("failed to open manifest file - %w", err)
	}
	jsonData, err := readAndUnseal(manifestFilepath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &manifest)
	if err != nil {
		// failed to unmarshal the json
		return nil, fmt.Errorf("failed to unmarshal manifest json - %w", err)
	}
	log.Info("Successfully loaded manifest from disk.")
	return &manifest, nil
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
		Issuer:       pkix.Name{CommonName: "obscuroCA"},
		Subject:      pkix.Name{CommonName: "obscuroUser"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
	}
	certPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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
	err = pem.Encode(certKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyOut,
	})
	return sealAndPersist(certKeyPEM.String(), userKeyFilepath)
}

func sealAndPersist(contents string, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s - %w", filepath, err)
	}
	defer f.Close()

	// START COMMENT-OUT IF DEBUGGING - while testing it can be useful to just f.WriteString(contents) below without sealing
	//    so that you can connect to edb using the persisted certs with mysql-client from the container.
	//    Note you also need to comment out the block in `readAndUnseal`

	// todo: do we prefer to seal with product key for upgradability or unique key to require fresh db with every code change
	enc, err := ecrypto.SealWithProductKey([]byte(contents), nil)
	if err != nil {
		return fmt.Errorf("failed to seal contents bytes with enclave key to persist in %s - %w", filepath, err)
	}
	// END COMMENT-OUT IF DEBUGGING

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

	// START COMMENT-OUT IF DEBUGGING - just `return b, nil` if debugging without sealing files

	data, err := ecrypto.Unseal(b, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to unseal data from file %s - %w", filepath, err)
	}
	return data, nil

	// END COMMENT-OUT IF DEBUGGING
}
