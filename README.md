# create-kalp-dapp

![Kalp DLT](/frontend/src/app/favicon.ico)

![Platform](https://img.shields.io/badge/platform-Kalp%20DLT-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![NPM Version](https://img.shields.io/npm/v/create-kalp-dapp)

A full-stack starter template featuring Next.js & Go, designed for building dApps and developing smart contracts on the Kalp DLT blockchain. This starter kit provides a simple greeting dApp to demonstrate interaction with a Kalp DLT smart contract.

![image](https://github.com/user-attachments/assets/d1d66943-d36e-4ba8-944b-6655ec159757)

## 🚀 Quick Start

```sh
npx create-kalp-dapp <your-dapp-name>

# cd into the directory
cd <your-dapp-name>
```

## 📦 Installation

Open your terminal and run:

```sh
npx create-kalp-dapp <your-dapp-name>
cd <your-dapp-name>
```

## 🏗 Project Structure

After creation, your project should look like this:

```
<your-dapp-name>/
├── backend/
│   ├── vendor/
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── krc.go
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   ├── hooks/
│   │   └── utils/
│   ├── .env.local.example
│   ├── next.config.mjs
│   ├── package.json
│   └── tailwind.config.ts
└── README.md
```

## 📜 Smart Contract Development

### Setting up the Smart Contract Project

1. Navigate to the smart contract directory:
   ```
   cd backend
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

### Smart Contract Code

The main smart contract file is `krc.go`. It includes a simple greeting contract with `Init`, `SetGreeting`, and `GetGreeting` functions.

### Compiling and Deploying the Smart Contract

1. Sign Up and Log In to [Kalp Studio Platform](https://console.kalp.studio/)
2. Go to Kalp Instant Deployer
3. Click on "Create New" Smart Contract
4. Enter the details:
   - Name: Enter a name for your smart contract
   - Category: Choose a category
   - Description: Add a description (optional)

   ![image](https://github.com/user-attachments/assets/39f41f16-a311-4427-8284-b9303872aa9e)

5. Upload your `backend.zip` file (zip the contents of the `backend` folder)

   ![image](https://github.com/user-attachments/assets/104f9955-05ce-4597-8348-628cf3e414ca)

6. Deploy your contract

After deployment, Kalp Studio will provide you with API endpoints for each function of your smart contract.

   ![image](https://github.com/user-attachments/assets/a223d7b9-d972-48b0-9aa2-06b7bea5a00f)

### Generating API Key

In Kalp Studio, navigate to the API key generation section and generate a new API key to authenticate your API requests.

## 💻 Frontend Development

### Setting up the Frontend

1. Navigate to the frontend directory:
   ```
   cd frontend
   ```

2. Install dependencies:
   ```
   npm install
   ```

3. Set up environment variables:
   - Copy `.env.local.example` to `.env.local`
   - Update `NEXT_PUBLIC_API_KEY` and `NEXT_PUBLIC_CONTRACT_ID` with your Kalp Studio API key and deployed contract ID

### Running the Frontend

Start the development server:

```
npm run dev
```

Visit `http://localhost:3000` to see your dApp in action.

![image](https://github.com/user-attachments/assets/afb392bd-4653-4325-a2d2-295d4527cac8)

## 🔧 Interacting with the Smart Contract

The frontend includes a simple interface to interact with your greeting smart contract:

- View the current greeting
- Set a new greeting

All interactions are handled through the Kalp DLT API, using the endpoints provided by Kalp Studio.

## 🔑 API Integration

The `useKalpApi` hook in `src/hooks/useKalpApi.tsx` handles all API calls to your smart contract. Make sure to update the contract ID and API key in your `.env.local` file for proper functionality.

## ➡️ Contributing

We welcome contributions! Please see our [CONTRIBUTING.md](https://github.com/Spydiecy/create-kalp-dapp/blob/main/CONTRIBUTING.md) file for details on how to get started.

## ⚖️ License

create-kalp-dapp is licensed under the [MIT License](https://github.com/Spydiecy/create-kalp-dapp/blob/main/LICENSE).


⭐️ If you find this project helpful, please give it a star on GitHub!
