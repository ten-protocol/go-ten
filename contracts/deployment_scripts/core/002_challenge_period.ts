import {HardhatRuntimeEnvironment} from 'hardhat/types';
import {DeployFunction} from 'hardhat-deploy/types';

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    const {deployments, getNamedAccounts} = hre;
    const {execute} = deployments;
    const {deployer} = await getNamedAccounts();

    // Get delay from environment variable, default to 0 for testnets
    const challengePeriod = process.env.CHALLENGE_PERIOD ?
        parseInt(process.env.CHALLENGE_PERIOD) : 0;

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
func.dependencies = ['ManagementContract']; // Ensure ManagementContract is deployed first