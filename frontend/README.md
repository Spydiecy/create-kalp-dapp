```ts
import React from 'react'

const page = () => {
  return (
    <div className='flex flex-col justify-center items-center'>

      <div className='border-2 py-4 px-16 mt-8 rounded-lg text-4xl w-fit'>Airdrop Machine</div>

      <div className='flex flex-col border-2 p-8 mt-8 rounded-lg w-fit gap-3'>
        Enter Your Address To Calim :
        <input placeholder='Enter your wallet address' type="text" className='border p-2 rounded-lg w-56 text-black' />
        <button className='border-2 p-2 rounded-lg bg-blue-500 hover:bg-blue-400 text-white'>Claim</button>
      </div>

      <div className='lg:flex gap-12'>
      <div className='flex flex-col border-2 p-8 mt-8 rounded-lg w-fit gap-3'>
        Total Airdrop Token Claimed :
        <p className='text-6xl font-bold w-56'>0</p>
      </div>

      <div className='flex flex-col border-2 p-8 mt-8 rounded-lg w-fit gap-3'>
        My Balance :
        <input placeholder='Enter your wallet address' type="text" className='border p-2 rounded-lg w-56 text-black' />
        <button className='border-2 p-2 rounded-lg bg-blue-500 hover:bg-blue-400 text-white'>See</button>

        <p className='text-2xl font-bold w-56'>Balance: <span className='text-blue-500 text-4xl'> {balance}</span></p>
      </div>

      </div>
      
    </div>
  )
}

export default page
```