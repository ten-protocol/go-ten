export default function getProgressPathAnimation(pathLength: number, progress: number) {
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
                    ease: 'linear',
                    repeatType: 'loop',
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
                    ease: 'easeInOut',
                },
                strokeDashoffset: {
                    duration: 1,
                    ease: 'easeInOut',
                },
                fillOpacity: {
                    delay: 1.5,
                    duration: 1,
                    ease: 'easeInOut',
                },
            },
        },
    };
}
