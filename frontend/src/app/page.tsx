"use client"
import React, { useEffect, useState } from 'react'
import { useKalpApi } from '@/hooks/useKalpAPI'

const Home: React.FC = () => {

 
    const { claim, balanceOf, totalSupply, loading } = useKalpApi();
    const [walletAddress, setWalletAddress] = useState("");
    const [balance, setBalance] = useState(0);
    const [totalAirdrop, setTotalAirdrop] = useState(0);
  
    const handleClaim = async () => {
      try {
        const data = await claim(walletAddress);
        await handleTotalSupply();
        console.log('Claim successful:', data);
      } catch (err) {
        console.error('Claim error:', err);
      }
    };
  
    const handleBalanceOf = async () => {
      try {
        const data = await balanceOf(walletAddress);
        setBalance(data.result.result)
        console.log('Balance:', data);
      } catch (err) {
        console.error('BalanceOf error:', err);
      }
    };
  
    const handleTotalSupply = async () => {
      try {
        const data = await totalSupply();
        setTotalAirdrop(data.result.result)
        console.log('Total Supply:', data);
      } catch (err) {
        console.error('TotalSupply error:', err);
      }
    };
  
    useEffect(() => {
      handleTotalSupply()
    }, []);
  
    return (
      <div className='flex flex-col justify-center items-center'>
  
        <div className='border-2 py-4 px-16 mt-8 rounded-lg text-4xl w-fit'>Airdrop Machine</div>
  
        <div className='flex flex-col border-2 p-8 mt-8 rounded-lg w-fit gap-3'>
          Enter Your Address To Calim :
          <input placeholder='Enter your wallet address' type="text" className='border p-2 rounded-lg w-56 text-black' onChange={(e) => setWalletAddress(e.target.value)} />
          <button className='border-2 p-2 rounded-lg bg-blue-500 hover:bg-blue-400 text-white disabled:bg-blue-400' onClick={handleClaim} disabled={loading}>{loading ? "Please wait.. " : "Claim"}</button>
        </div>
  
        <div className='lg:flex gap-12'>
          <div className='flex flex-col border-2 p-8 mt-8 rounded-lg w-fit gap-3'>
            Total Airdrop Token Claimed :
            <p className='text-6xl text-blue-500 font-bold w-56'>{totalAirdrop}</p>
          </div>
  
          <div className='flex flex-col border-2 p-8 mt-8 rounded-lg w-fit gap-3'>
            My Balance :
            <input placeholder='Enter your wallet address' type="text" className='border p-2 rounded-lg w-56 text-black' onChange={(e) => setWalletAddress(e.target.value)} />
            <button className='border-2 p-2 rounded-lg bg-blue-500 hover:bg-blue-400 text-white' onClick={handleBalanceOf}>See</button>
  
            <p className='text-2xl font-bold w-56'>Balance: <span className='text-blue-500 text-4xl'> {balance}</span></p>
          </div>
  
        </div>
  
      </div>
    )
  }

export default Home;