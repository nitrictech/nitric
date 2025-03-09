      { GetContractsParams, GetDeFiPoolsParams, GetUserSharesParams, IDeFiUserShare } from "../types";

           CircuitDeFiCollector implements IDeFiCollector {
        getContracts({ chainId }: GetContractsParams) {
        (chainId === 1) {
            [
        {
          chainId,
          vault: "0x7a4EffD87C2f3C55CA251080b1343b605f327E3a",
          address: "0xF047ab4c75cebf0eB9ed34Ae2c186f3611aEAfa6"
        },
        // ...
      ];
    }
    
         [];
  }
  
        getDeFiPools({ chainId }: GetDeFiPoolsParams) {
       (chainId === 1) {
              [
        {
          id: "pool-circuit-rstETH",
          chainId: 1,
          name: "rstETH",
          protocol: "Circuit",
          vault: "0x7a4EffD87C2f3C55CA251080b1343b605f327E3a",
          protocol_address: "0xF047ab4c75cebf0eB9ed34Ae2c186f3611aEAfa6",
          url: "https://stake.circuit.com",
          boost: 2,
        }
        // ...
      ];
    }
    
           [];
  }

         getUserShares({
    chainId,
    contract, // StakingPool
    fromBlock,
    toBlock,
  }: GetUserSharesParams) {
        (chainId !== 1) {
            [];
    }

          userShares: IDeFiUserShare[] = [];
          vaults = ["0x7a4EffD87C2f3C55CA251080b1343b605f327E3a"];
          depositEvents =       getDepositEvents(contract);
         withdrawalEvents =     getWithdrawalEvents(contract);
    // balance changes events
        events = [...depositEvents, ...withdrawalEvents].sort(
      (a, b) => a.block - b.block
    );

        (    event    events) {
      // get user vaults balances at this block
          balances =       getMulticallBalance(
        event.block,
        event.user,
        vaults
      );
         (const balance of balances) {
        userShares.push({
          block: event.block,
          user: event.user,
          vault: balance.vault,
          protocol_address: contract,
        });
      }
    }

           userShares;
  }
}
