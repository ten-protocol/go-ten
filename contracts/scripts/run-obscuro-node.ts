import { subtask, task } from "hardhat/config";
import { TASK_COMPILE } from "hardhat/builtin-tasks/task-names";
import * as path from "path";
import { HardhatPluginError } from 'hardhat/plugins';

import * as dockerApi from 'node-docker-api';
import YAML from 'yaml'
import * as fs from 'fs';


task("run-obscuro-node", "Starts up the obscuro node docker container.")
.addParam('composeFile', 'Location to docker compose', "./scripts/docker-compose.non-sgx.yml")
.setAction(async function(args, hre) {
    const docker = new dockerApi.Docker({ socketPath: '/var/run/docker.sock' })

    const file = fs.readFileSync(args.composeFile, 'utf8');
    const composeConfig = YAML.parse(file);

    const networkName = composeConfig.networks.default.name;

    const networks = await docker.network.list()
    const networkNames = networks.map((network)=> (network.data as any).Name);

    const networkToUse = networkNames.find((x)=>x == networkName);
    const hasNetwork = networkToUse != null;

    console.log(`Network[${networkName}].exists=${hasNetwork}`);

    if (!hasNetwork) {
        const result = await docker.network.create({
            Name: networkName
        })
        console.log(`Network[${result.id}] created.`);
    }


    async function spawnContainer(config: any) {
        console.log(`env - ${JSON.stringify(config.environment)} type of ${typeof(config.environment)}`);
        const envMap = new Map(Object.entries(config.environment));
        envMap.forEach((v,k)=> {
            console.log(v,k);
        })



        return;
        docker.container.create({
            "Hostname": "",
            "Domainname": "",
            "User": "",
            "AttachStdin": false,
            "AttachStdout": true,
            "AttachStderr": true,
            "Tty": false,
            "OpenStdin": false,
            "StdinOnce": false,
            Env: config.environment,
            "Entrypoint": "",
            Image: config.image,
            "Volumes": {
              "/volumes/data": {}
            },
            "Healthcheck":{
               "Test": ["CMD-SHELL", "curl localhost:3000"],
               "Interval": 1000000000,
               "Timeout": 10000000000,
               "Retries": 10,
               "StartPeriod": 60000000000
            },
            "WorkingDir": "",
            "NetworkDisabled": false,
            "MacAddress": "12:34:56:78:9a:bc",
            "ExposedPorts": {
                    "22/tcp": {}
            },
            "StopSignal": "SIGTERM",
            "HostConfig": {
              "Binds": ["/tmp:/tmp"],
              "Tmpfs": { "/run": "rw,noexec,nosuid,size=65536k" },
              "Links": ["redis3:redis"],
              "Memory": 0,
              "MemorySwap": 0,
              "MemoryReservation": 0,
              "KernelMemory": 0,
              "CpuPercent": 80,
              "CpuShares": 512,
              "CpuPeriod": 100000,
              "CpuQuota": 50000,
              "CpusetCpus": "0,1",
              "CpusetMems": "0,1",
              "IOMaximumBandwidth": 0,
              "IOMaximumIOps": 0,
              "BlkioWeight": 300,
              "BlkioWeightDevice": [{}],
              "BlkioDeviceReadBps": [{}],
              "BlkioDeviceReadIOps": [{}],
              "BlkioDeviceWriteBps": [{}],
              "BlkioDeviceWriteIOps": [{}],
              "MemorySwappiness": 60,
              "OomKillDisable": false,
              "OomScoreAdj": 500,
              "PidMode": "",
              "PidsLimit": -1,
              "PortBindings": { "22/tcp": [{ "HostPort": "11022" }] },
              "PublishAllPorts": false,
              "Privileged": false,
              "ReadonlyRootfs": false,
              "Dns": ["8.8.8.8"],
              "DnsOptions": [""],
              "DnsSearch": [""],
              "ExtraHosts": null,
              "VolumesFrom": ["parent", "other:ro"],
              "CapAdd": ["NET_ADMIN"],
              "CapDrop": ["MKNOD"],
              "GroupAdd": ["newgroup"],
              "RestartPolicy": { "Name": "", "MaximumRetryCount": 0 },
              "NetworkMode": "bridge",
              "Devices": [],
              "Sysctls": { "net.ipv4.ip_forward": "1" },
              "Ulimits": [{}],
              "LogConfig": { "Type": "json-file", "Config": {} },
              "SecurityOpt": [],
              "StorageOpt": {},
              "CgroupParent": "",
              "VolumeDriver": "",
              "ShmSize": 67108864
           },
           "NetworkingConfig": {
               "EndpointsConfig": {
                   "isolated_nw" : {
                       "IPAMConfig": {
                           "IPv4Address":"172.20.30.33",
                           "IPv6Address":"2001:db8:abcd::3033",
                           "LinkLocalIPs":["169.254.34.68", "fe80::3468"]
                       },
                       "Links":["container_1", "container_2"],
                       "Aliases":["server_x", "server_y"]
                   }
               }
           }
        })
    }

    const hostConfig = composeConfig.services.host
    spawnContainer(hostConfig);
    
    const enclaveConfig = composeConfig.services.encalve


});