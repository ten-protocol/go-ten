'use client';

import { useEffect, useState, useRef } from 'react';
import { useInView } from 'framer-motion';

type Props = {
    text: string;
    hover?: boolean;
    active?: boolean;
    onView?: boolean;
    glitchDuration?:number
};

export default function GlitchTextAnimation({
    text,
    hover = true,
    active = true,
    onView = true,
    glitchDuration = 400
}: Props) {
    const [displayText, setDisplayText] = useState(text);
    const [isGlitching, setIsGlitching] = useState(false);
    const [glitchVariants, setGlitchVariants] = useState({
        redOffset: 0,
        cyanOffset: 0,
        duration: 0.3,
        redHue: 0,
        cyanHue: 0,
        redSkew: 0,
        cyanSkew: 0,
    });
    const glitchTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const ref = useRef<HTMLParagraphElement>(null);
    const isInView = useInView(ref, { once: false });

    const generateRandomGlitchVariants = () => {
        return {
            redOffset: Math.floor(Math.random() * 2) + 0.5, // 0.5-2.5px
            cyanOffset: Math.floor(Math.random() * 2) + 0.5, // 0.5-2.5px
            duration: (Math.random() * 0.15) + 0.2, // 0.2-0.35s
            redHue: Math.floor(Math.random() * 20) - 10, // -10 to +10 hue shift
            cyanHue: Math.floor(Math.random() * 20) - 10, // -10 to +10 hue shift
            redSkew: (Math.random() * 20) - 10, // -10 to +10 degrees
            cyanSkew: (Math.random() * 20) - 10, // -10 to +10 degrees
        };
    };

    const triggerGlitch = () => {
        if (isGlitching) return;
        
        const newVariants = generateRandomGlitchVariants();
        setGlitchVariants(newVariants);
        
        setIsGlitching(true);
        
        if (glitchTimeoutRef.current) {
            clearTimeout(glitchTimeoutRef.current);
        }

        glitchTimeoutRef.current = setTimeout(() => {
            setIsGlitching(false);
        }, glitchDuration);
    };

    useEffect(() => {
        if (active && onView && isInView) {
            triggerGlitch();
        } else if (active && !onView) {
            triggerGlitch();
        }
    }, [isInView, active, onView]);

    useEffect(() => {
        return () => {
            if (glitchTimeoutRef.current) {
                clearTimeout(glitchTimeoutRef.current);
            }
        };
    }, []);

    const handleInteraction = () => {
        if (hover && !isGlitching) {
            triggerGlitch();
        }
    };

    return (
        <span className="relative inline-block">
            <span 
                ref={ref} 
                onMouseOver={handleInteraction} 
                onFocus={handleInteraction}
                className={`transition-all duration-75 relative z-20 ${
                    isGlitching 
                        ? 'text-white' 
                        : ''
                }`}
                style={{
                    animation: isGlitching 
                        ? 'glitch-shake 0.3s ease-in-out' 
                        : 'none',
                    transform: isGlitching 
                        ? `skew(${glitchVariants.redSkew * 0.5}deg)` 
                        : 'none'
                }}
            >
                {displayText}
            </span>
            
            {isGlitching && (
                <>
                    <span
                        className="absolute inset-0 opacity-60 z-10"
                        style={{
                            animation: `glitch-red-1 ${glitchVariants.duration}s ease-in-out`,
                            transform: `translate(${glitchVariants.redOffset}px, 0) skew(${glitchVariants.redSkew}deg)`,
                            filter: `hue-rotate(${glitchVariants.redHue}deg)`,
                            color: 'hsl(0, 100%, 60%)',
                        }}
                    >
                        {displayText}
                    </span>
                    <span 
                        className="absolute inset-0 opacity-40 z-10"
                        style={{
                            animation: `glitch-red-2 ${glitchVariants.duration}s ease-in-out`,
                            transform: `translate(${glitchVariants.redOffset - 0.5}px, 0.5px) skew(${glitchVariants.redSkew + 3}deg)`,
                            filter: `hue-rotate(${glitchVariants.redHue + 5}deg)`,
                            color: 'hsl(0, 100%, 50%)',
                        }}
                    >
                        {displayText}
                    </span>
                    
                    <span
                        className="absolute inset-0 opacity-60 z-10"
                        style={{
                            animation: `glitch-cyan-1 ${glitchVariants.duration}s ease-in-out`,
                            transform: `translate(-${glitchVariants.cyanOffset}px, 0) skew(${glitchVariants.cyanSkew}deg)`,
                            filter: `hue-rotate(${glitchVariants.cyanHue}deg)`,
                            color: 'hsl(180, 100%, 60%)',
                        }}
                    >
                        {displayText}
                    </span>
                    <span 
                        className="absolute inset-0 opacity-40 z-10"
                        style={{
                            animation: `glitch-cyan-2 ${glitchVariants.duration}s ease-in-out`,
                            transform: `translate(-${glitchVariants.cyanOffset - 0.5}px, -0.5px) skew(${glitchVariants.cyanSkew + 3}deg)`,
                            filter: `hue-rotate(${glitchVariants.cyanHue + 5}deg)`,
                            color: 'hsl(180, 100%, 50%)',
                        }}
                    >
                        {displayText}
                    </span>
                    <span
                        className="absolute inset-0 opacity-35 z-10"
                        style={{
                            animation: `glitch-artifact-1 ${glitchVariants.duration}s ease-in-out`,
                            transform: `translate(${glitchVariants.redOffset + 0.5}px, -${glitchVariants.redOffset}px) skew(${glitchVariants.redSkew + 5}deg)`,
                            filter: `hue-rotate(${Math.random() * 360}deg)`,
                            color: 'hsl(0, 0%, 80%)',
                        }}
                    >
                        {displayText}
                    </span>
                    <span 
                        className="absolute inset-0 opacity-30 z-10"
                        style={{
                            animation: `glitch-artifact-2 ${glitchVariants.duration}s ease-in-out`,
                            transform: `translate(-${glitchVariants.cyanOffset + 0.5}px, ${glitchVariants.cyanOffset}px) skew(${glitchVariants.cyanSkew + 5}deg)`,
                            filter: `hue-rotate(${Math.random() * 360}deg)`,
                            color: 'hsl(60, 100%, 60%)',
                        }}
                    >
                        {displayText}
                    </span>
                </>
            )}
        </span>
    );
}
