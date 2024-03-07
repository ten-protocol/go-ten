import { task } from "hardhat/config";
import http from 'http';


task("obscuro:wallet-extension:add-key", "Creates a viewing key for a specifiec address")
.addParam("address", "The address for which to add key")
.setAction(async function(args, hre) {
    async function viewingKeyForAddress(address: string) : Promise<string> {
        return new Promise((resolve, fail)=> {
    
            const data = {"address": address}
    
            const req = http.request({
                host: '127.0.0.1',
                port: 3000,
                path: '/generateviewingkey/',
                method: 'post',
                headers: {
                    'Content-Type': 'application/json'
                }
            }, (response)=>{
                if (response.statusCode != 200) {
                    console.error(response);
                    fail(response.statusCode);
                    return;
                }
    
                let chunks : string[] = []
                response.on('data', (chunk)=>{
                    chunks.push(chunk);
                });
    
                response.on('end', ()=> { 
                    resolve(chunks.join('')); 
                });
            });
            req.write(JSON.stringify(data));
            req.end()
            setTimeout(resolve, 15_000);
        });
    }
    
    interface SignedData { signature: string, address: string }
    
    async function submitKey(signedData: SignedData) : Promise<number> {
        return await new Promise(async (resolve, fail)=>{ 
            const req = http.request({
                host: '127.0.0.1',
                port: 3000,
                path: '/submitviewingkey/',
                method: 'post',
                headers: {
                    'Content-Type': 'application/json'
                }
            }, (response)=>{
                if (response.statusCode == 200) { 
                    resolve(response.statusCode);
                } else {
                    console.log(response.statusMessage)
                    fail(response.statusCode);
                }
            });
    
            req.write(JSON.stringify(signedData));
            req.end()
        });
    }

    const key = await viewingKeyForAddress(args.address);

    const signaturePromise = (await hre.ethers.getSigner(args.address)).signMessage(`vk${key}`);
    const signedData = { 'signature': await signaturePromise, 'address': args.address };
    await submitKey(signedData)
});