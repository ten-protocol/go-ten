'use client';

import { ReactNode, useRef } from 'react';

import { useInView } from 'framer-motion';

type Props = {
    children: ReactNode;
};

export default function PrimaryCard({ children }: Props) {
    const ref = useRef<HTMLDivElement>(null);
    const isInView = useInView(ref, { margin: '0px 100px -50px 0px', once: true });

    return (
        <div className="relative h-full">
            <div
                ref={ref}
                className="shape-one flex items-center relative backdrop-blur-xl h-full"
                style={{ transformStyle: 'preserve-3d' }}
            >
                <div
                    className="z-10 animate-scan-overlay w-full absolute pointer-events-none"
                    style={{
                        height: isInView ? '100%' : '0',
                        transition: 'all .3s cubic-bezier(0.17, 0.55, 0.55, 1) 0.3s',
                    }}
                />

                <div
                    className="z-0 m-8 grow"
                    style={{
                        opacity: isInView ? 1 : 0,
                        filter: isInView ? 'blur(0)' : 'blur(3rem)',
                        transition: 'all .3s cubic-bezier(0.17, 0.55, 0.55, 1) .6s',
                    }}
                >
                    {children}
                </div>
            </div>
        </div>
    );
}
