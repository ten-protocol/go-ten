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
});