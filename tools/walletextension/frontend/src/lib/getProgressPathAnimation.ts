import { Variants } from 'framer-motion';

export default function getProgressPathAnimation(pathLength: number, progress: number): Variants {
    const visibleDashLength = pathLength * (progress / 100);
    const gapLength = pathLength - visibleDashLength;
    const pathAnimationDuration = pathLength / 20;

    return {
        initial: {
            strokeDashoffset: pathLength,
            fillOpacity: 0,
            strokeDasharray: `${visibleDashLength} ${gapLength}`,
        },
        animate: {
            strokeDashoffset: [pathLength, 0],
            strokeDasharray: `${visibleDashLength} ${gapLength}`,
            transition: {
                strokeDashoffset: {
                    duration: pathAnimationDuration,
                    repeat: Infinity,
                    ease: [0.4, 0, 0.2, 1],
                    repeatType: "loop" as const,
                },
            },
        },
        complete: {
            strokeDashoffset: 0,
            strokeDasharray: `${pathLength} 0`,
            fillOpacity: 1,
            transition: {
                strokeDasharray: {
                    duration: 1,
                    ease: [0.4, 0, 0.2, 1],
                },
                strokeDashoffset: {
                    duration: 1,
                    ease: [0.4, 0, 0.2, 1],
                },
                fillOpacity: {
                    delay: 1.5,
                    duration: 1,
                    ease: [0.4, 0, 0.2, 1],
                },
            },
        },
        error: {
            strokeDashoffset: 0,
            strokeDasharray: `0 ${pathLength}`,
            opacity: 0,
            transition: {
                duration: 1,
                ease: [0.4, 0, 0.2, 1],
            },
        },
    };
}
