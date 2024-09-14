export interface ApiResponse {
    data: string;
    error?: string;
  }
  
  export const initializeContract = async (
    name: string,
    symbol: string,
    decimals: string
  ): Promise<ApiResponse> => {
    const response = await fetch('/api/initialize', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name, symbol, decimals }),
    });
    const data = await response.json();
    return data;
  };
  
  export const mintTokens = async (
    amount: string
  ): Promise<ApiResponse> => {
    const response = await fetch('/api/mint', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ amount }),
    });
    const data = await response.json();
    return data;
  };
  
  export const transferTokens = async (
    recipient: string,
    amount: string
  ): Promise<ApiResponse> => {
    const response = await fetch('/api/transfer', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ recipient, amount }),
    });
    const data = await response.json();
    return data;
  };
  
  export const getTotalSupply = async ()  => {
    const response = await fetch('/api/totalSupply');
    const data = await response.json();
    return data;
  };