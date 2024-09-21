'use client'

import { useState, useEffect } from 'react'
import { useKalpApi } from '../hooks/useKalpAPI'

export default function Home() {
  const [greeting, setGreeting] = useState('')
  const [newGreeting, setNewGreeting] = useState('')
  const { getGreeting, setGreeting: updateGreeting, loading, error } = useKalpApi()

  useEffect(() => {
    handleGetGreeting()
  }, [])

  const handleGetGreeting = async () => {
    try {
      const response = await getGreeting()
      console.log('Get Greeting Response:', response);
      if (response && response.result) {
        console.log('Result data:', response.result);
        
        // Parse the result if it's a string
        if (typeof response.result === 'string') {
          try {
            const parsedResult = JSON.parse(response.result);
            if (parsedResult && parsedResult.result) {
              setGreeting(parsedResult.result);
            } else {
              setGreeting(response.result); // Fallback to the original string if parsing fails
            }
          } catch (parseError) {
            console.error("Error parsing result:", parseError);
            setGreeting(response.result); // Use the original string if parsing fails
          }
        } else if (typeof response.result === 'object' && response.result.result) {
          setGreeting(response.result.result);
        } else {
          setGreeting("Unexpected result format");
        }
      } else {
        setGreeting("No greeting data found");
      }
    } catch (err) {
      console.error("Failed to get greeting", err)
      setGreeting("Error fetching greeting")
    }
  }

  const handleSetGreeting = async () => {
    try {
      const response = await updateGreeting(newGreeting)
      console.log('Set Greeting Response:', response);
      setNewGreeting('')
      handleGetGreeting()
    } catch (err) {
      console.error("Failed to set greeting", err)
    }
  }

  return (
    <div className="max-w-2xl mx-auto">
      <h2 className="text-2xl font-bold mb-8 text-center">Interact with the Greeting Contract</h2>
      
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
          <strong className="font-bold">Error:</strong>
          <span className="block sm:inline"> {error.message}</span>
        </div>
      )}

      <div className="bg-white shadow-md rounded-lg p-6 mb-8">
        <h3 className="text-xl font-semibold mb-4">Current Greeting</h3>
        <div className="flex items-center justify-between">
          <p className="text-lg">
            {greeting}
          </p>
          <button
            onClick={handleGetGreeting}
            className="btn btn-primary"
            disabled={loading}
          >
            {loading ? 'Loading...' : 'Refresh Greeting'}
          </button>
        </div>
      </div>

      <div className="bg-white shadow-md rounded-lg p-6">
        <h3 className="text-xl font-semibold mb-4">Set New Greeting</h3>
        <div className="flex items-center">
          <input
            type="text"
            value={newGreeting}
            onChange={(e) => setNewGreeting(e.target.value)}
            placeholder="Enter new greeting"
            className="input flex-grow mr-2"
          />
          <button
            onClick={handleSetGreeting}
            className="btn btn-secondary whitespace-nowrap"
            disabled={loading}
          >
            {loading ? 'Setting...' : 'Set Greeting'}
          </button>
        </div>
      </div>
    </div>
  )
}