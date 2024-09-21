import "./globals.css";

export const metadata = {
  title: 'Kalp DLT Greeting dApp',
  description: 'A simple greeting dApp built on Kalp DLT',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <link rel="icon" href="/favicon.ico" sizes="any" />
      </head>
      <body className="flex flex-col min-h-screen">
        <header className="bg-black text-white py-6">
          <div className="container mx-auto px-4">
            <div className="flex flex-col items-center justify-center">
              <img 
                src="/favicon.ico" 
                alt="Kalp DLT Logo" 
                className="w-16 h-16 mb-2" 
                aria-label="Kalp DLT Logo"
              />
              <h1 className="text-2xl font-bold text-center">Kalp DLT Greeting dApp</h1>
            </div>
          </div>
        </header>
        <main className="container mx-auto px-4 py-8 flex-grow">
          {children}
        </main>
        <footer className="bg-black text-white py-4 mt-8">
          <div className="container mx-auto px-4 text-center">
            <p className="text-sm">Built with create-kalp-dapp | © 2024 Kalp DLT</p>
          </div>
        </footer>
      </body>
    </html>
  );
}