'use client';

import { useState } from 'react';
import {
  initializeContract,
  mintTokens,
  transferTokens,
  getTotalSupply,
} from '@/app/lib/api';


const Home: React.FC = () => {
  const [name, setName] = useState('');
  const [symbol, setSymbol] = useState('');
  const [decimals, setDecimals] = useState('');
  const [mintAmount, setMintAmount] = useState('');
  const [recipient, setRecipient] = useState('');
  const [transferAmount, setTransferAmount] = useState('');
  const [totalSupply, setTotalSupply] = useState(null);
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState<'success' | 'error'>('success');

  


  // Handle contract initialization
  const handleInitialize = async () => {
    try {
      const result = await initializeContract(name, symbol, decimals);
      setMessage(`Contract initialized successfully.`);
      setMessageType('success');
      console.log('Initialize Result:', result);
    } catch (error) {
      setMessage('Failed to initialize contract.');
      setMessageType('error');
      console.error('Initialize Error:', error);
    }
  };

  // Handle minting tokens
  const handleMint = async () => {
    try {
      const result = await mintTokens(mintAmount);
      setMessage(`Minted ${mintAmount} tokens successfully.`);
      setMessageType('success');
      console.log('Mint Result:', result);
    } catch (error) {
      setMessage('Failed to mint tokens.');
      setMessageType('error');
      console.error('Mint Error:', error);
    }
  };

  // Handle transferring tokens
  const handleTransfer = async () => {
    try {
      const result = await transferTokens(recipient, transferAmount);
      setMessage(`Transferred ${transferAmount} tokens to ${recipient}.`);
      setMessageType('success');
      console.log('Transfer Result:', result);
    } catch (error) {
      setMessage('Failed to transfer tokens.');
      setMessageType('error');
      console.error('Transfer Error:', error);
    }
  };

  // Handle getting total supply
  const handleGetTotalSupply = async () => {
    try {
      const result = await getTotalSupply();
      setTotalSupply(result);
      setMessage('Fetched total supply successfully.');
      setMessageType('success');
      console.log('Total Supply:', result);
    } catch (error) {
      setMessage('Failed to fetch total supply.');
      setMessageType('error');
      console.error('Total Supply Error:', error);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <h1 className="text-3xl font-bold mb-8 text-center">
        Hyperledger ERC-20 Token Interface
      </h1>

      {/* Message Display */}
      {message && (
        <div
          className={`${
            messageType === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
          } p-4 mb-6 rounded`}
        >
          {message}
        </div>
      )}

      {/* Initialize Contract Section */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-4">Initialize Contract</h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <input
            type="text"
            placeholder="Token Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="input"
          />
          <input
            type="text"
            placeholder="Token Symbol"
            value={symbol}
            onChange={(e) => setSymbol(e.target.value)}
            className="input"
          />
          <input
            type="text"
            placeholder="Decimals"
            value={decimals}
            onChange={(e) => setDecimals(e.target.value)}
            className="input"
          />
          <button onClick={handleInitialize} className="btn-primary">
            Initialize
          </button>
        </div>
      </section>

      {/* Mint Tokens Section */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-4">Mint Tokens</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-end">
          <input
            type="number"
            placeholder="Amount to Mint"
            value={mintAmount}
            onChange={(e) => setMintAmount(e.target.value)}
            className="input"
          />
          <div className="hidden md:block"></div> {/* Spacer */}
          <button onClick={handleMint} className="btn-primary">
            Mint Tokens
          </button>
        </div>
      </section>

      {/* Transfer Tokens Section */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-4">Transfer Tokens</h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <input
            type="text"
            placeholder="Recipient Address"
            value={recipient}
            onChange={(e) => setRecipient(e.target.value)}
            className="input col-span-2"
          />
          <input
            type="number"
            placeholder="Amount to Transfer"
            value={transferAmount}
            onChange={(e) => setTransferAmount(e.target.value)}
            className="input"
          />
          <button onClick={handleTransfer} className="btn-primary">
            Transfer
          </button>
        </div>
      </section>

      {/* Total Supply Section */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-4">Total Supply</h2>
        <button onClick={handleGetTotalSupply} className="btn-primary mb-4">
          Get Total Supply
        </button>
        {totalSupply && (
          <pre className="bg-gray-100 p-4 rounded overflow-x-auto">
            {JSON.stringify(totalSupply, null, 2)}
          </pre>
        )}
      </section>
    </div>
  );
};

export default Home;