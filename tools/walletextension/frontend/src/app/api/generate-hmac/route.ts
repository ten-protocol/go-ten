import { NextRequest, NextResponse } from 'next/server';
import { HMAC_SECRET } from '@/lib/constants';
import crypto from 'crypto';

export async function POST(request: NextRequest) {
    try {
        const body = await request.json();
        const { timestamp } = body;

        // Validate required fields
        if (!timestamp) {
            return NextResponse.json(
                { error: 'Missing required field: timestamp' },
                { status: 400 }
            );
        }

        // Check if HMAC secret is configured
        if (!HMAC_SECRET) {
            return NextResponse.json(
                { error: 'HMAC secret not configured' },
                { status: 500 }
            );
        }

        // Generate HMAC signature
        const hmac = crypto.createHmac('sha256', HMAC_SECRET);
        hmac.update(timestamp);
        const signature = hmac.digest('hex');

        return NextResponse.json({ signature });
    } catch (error) {
        console.error('Error in generate-hmac API:', error);
        return NextResponse.json(
            { error: 'Internal server error' },
            { status: 500 }
        );
    }
} 