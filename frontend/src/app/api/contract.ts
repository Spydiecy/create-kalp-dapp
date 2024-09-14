import type { NextApiRequest, NextApiResponse } from 'next';
import { callKalpApi } from '@/utils/apiHelper';

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method === 'POST') {
    const { name, symbol, decimals } = req.body;
    try {
      const { status, data } = await callKalpApi(
        'https://gateway-api.kalp.studio/v1/contract/kalp/invoke/Ymx4MuFOcZqP2PEizNu5Yl04FPVstT7f1726138235031/Initialize',
        { name, symbol, decimals }
      );
      res.status(status).json(data);
    } catch (error) {
      res
        .status(500)
        .json({ error: 'Initialization failed', details: error });
    }
  } else {
    res.status(405).json({ message: 'Method not allowed' });
  }
}