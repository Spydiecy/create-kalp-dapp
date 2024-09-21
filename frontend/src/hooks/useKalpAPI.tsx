"use client"

import { useState } from 'react';

export const useKalpApi = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const apiKey = process.env.NEXT_PUBLIC_API_KEY;
  const contractId = process.env.NEXT_PUBLIC_CONTRACT_ID;

  const callApi = async (endpoint: string, args: { [key: string]: any } = {}) => {
    if (!contractId) {
      throw new Error('Contract ID is not set. Please check your environment variables.');
    }

    setError(null);
    setLoading(true);

    const params = {
      network: 'TESTNET',
      blockchain: 'KALP',
      walletAddress: '', // This should be replaced with the actual user's wallet address in a real application
      args: args,
    };

    try {
      console.log(`Calling API: ${endpoint}`, params);
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'x-api-key': apiKey || '',
        },
        body: JSON.stringify(params),
      });

      const data = await response.json();
      console.log(`Full API Response:`, data);
      
      if (data.result) {
        console.log(`Result data:`, data.result);
      } else {
        console.log(`No result data found in the response`);
      }

      if (!response.ok) {
        throw new Error(data.message || `API call failed with status ${response.status}`);
      }

      setLoading(false);
      return data;
    } catch (err: any) {
      console.error(`API Error:`, err);
      setError(err);
      setLoading(false);
      throw err;
    }
  };

  const getGreeting = async () => {
    const endpoint = `https://gateway-api.kalp.studio/v1/contract/kalp/query/${contractId}/GetGreeting`;
    return callApi(endpoint);
  };

  const setGreeting = async (greeting: string) => {
    const endpoint = `https://gateway-api.kalp.studio/v1/contract/kalp/invoke/${contractId}/SetGreeting`;
    const args = { greeting };
    return callApi(endpoint, args);
  };

  return { getGreeting, setGreeting, loading, error };
};