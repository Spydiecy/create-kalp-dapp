"use client"
import { useState } from 'react';

export const useKalpApi = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const apiKey = process.env.NEXT_PUBLIC_API_KEY;

  const callApi = async (endpoint: string, args : { [key: string]: any }) => {
    setError(null);
    const params = {
      network: 'TESTNET',
      blockchain: 'KALP',
      walletAddress: '928bc86952ebb55788e2042ad478b8c1db3ded0d',
      args: args,
    };

    try {
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'x-api-key': apiKey!,
        },
        body: JSON.stringify(params),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Something went wrong');
      }
      setLoading(false);
      return data;
    } catch (err : any) {
      setError(err);
      setLoading(false);
      throw err;
    }
  };

  const claim = async (address : string) => {
    setLoading(true);
    const endpoint =
      'https://gateway-api.kalp.studio/v1/contract/kalp/invoke/M5fb2tJjsQk6unfWDttcTxkp1GC129dr1726253658000/Claim';
    const args = {
      amount: 100,
      address: address,
    };
    return callApi(endpoint, args);
  };

  const balanceOf = async (account : string) => {
    const endpoint =
      'https://gateway-api.kalp.studio/v1/contract/kalp/query/M5fb2tJjsQk6unfWDttcTxkp1GC129dr1726253658000/BalanceOf';
    const args = {
      account: account,
    };
    return callApi(endpoint, args);
  };

  const totalSupply = async () => {
    const endpoint =
      'https://gateway-api.kalp.studio/v1/contract/kalp/query/M5fb2tJjsQk6unfWDttcTxkp1GC129dr1726253658000/TotalSupply';
    const args = {};
    return callApi(endpoint, args);
  };

  return { claim, balanceOf, totalSupply, loading, error };
};


