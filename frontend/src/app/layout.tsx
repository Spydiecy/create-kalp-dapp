import "./globals.css";
import { Handjet } from 'next/font/google';

const inter = Handjet({ subsets: ['latin'] });



export const metadata = {
  title: 'Airdrop Machine',
  description: 'Click to get airdrop',
};



export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={inter.className}>{children}</body>
    </html>
  );
}
