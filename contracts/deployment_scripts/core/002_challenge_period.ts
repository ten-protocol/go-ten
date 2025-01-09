import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployments, getNamedAccounts} = hre;
    const {execute} = deployments;
    const {deployer} = await getNamedAccounts();

    const challengePeriod = process.env.L1_CHALLENGE_PERIOD ? 
        Number(process.env.L1_CHALLENGE_PERIOD) : 0;

    console.log(`Setting challenge period to ${challengePeriod}`);

    await execute(
        'ManagementContract',
        {from: deployer, log: true},
        'setChallengePeriod',
        challengePeriod
    );
};

export default func;
func.tags = ['ChallengePeriod'];
func.dependencies = ['ManagementContract'];