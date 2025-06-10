import { motion } from 'framer-motion';
import getProgressPathAnimation from '@/lib/getProgressPathAnimation';

type Props = {
    progress?: number;
    error?: boolean;
};

export default function ConnectWalletProgressAnimation({ progress = 0, error = true }: Props) {
    const circleRadius = 141.277;
    const circleCircumference = 2 * Math.PI * circleRadius;
    const progressLength = (progress / 100) * circleCircumference;
    const complete = progress >= 100;

    const mainCircleVariants = {
        initial: {
            rotate: 0,
            strokeDasharray: `0 ${circleCircumference}`,
            stroke: '#e3e3e3',
        },
        animate: {
            rotate: 360,
            strokeDasharray: `${progressLength} ${circleCircumference - progressLength}`,
            stroke: '#e3e3e3',
            transition: {
                rotate: {
                    duration: 3,
                    repeat: Infinity,
                    ease: 'linear',
                },
                strokeDasharray: {
                    duration: 0.5,
                    ease: 'easeOut',
                },
            },
        },
        complete: {
            rotate: 0,
            strokeDasharray: `${circleCircumference} 0`,
            stroke: '#00c951',
            transition: {
                duration: 1,
                ease: 'easeInOut',
            },
        },
        error: {
            strokeDasharray: `0 ${circleCircumference}`,
            opacity: 0,
            transition: {
                duration: 0.5,
                ease: 'easeInOut',
            },
        },
    };

    const backgroundCircleVariants1 = {
        initial: {
            stroke: 'rgba(255,255,255,0.3)',
        },
        animate: {
            stroke: 'rgba(255,255,255,0.3)',
        },
        complete: {
            rotate: 0,
            x: 0,
            y: 0,
            stroke: 'rgba(76,175,80,0.3)',
            transition: {
                duration: 1,
                ease: 'easeInOut',
            },
        },
        error: {
            rotate: 0,
            x: 0,
            y: 0,
            stroke: 'rgba(255,103,103,0.8)',
            transition: {
                duration: 1,
                ease: 'easeInOut',
            },
        },
    };

    const variant = error ? 'error' : complete ? 'complete' : 'animate';

    const backgroundPath = {
        fill: 'none',
        stroke: error ? 'rgba(255,132,132,0.4)' : 'rgba(227, 227, 227, 0.25)',
        strokeWidth: '1px',
    };

    return (
        <div>
            <svg
                width="100%"
                height="100%"
                viewBox="0 0 300 300"
                style={{
                    fillRule: 'evenodd',
                    clipRule: 'evenodd',
                    strokeLinecap: 'round',
                    strokeLinejoin: 'round',
                    strokeMiterlimit: 1.5,
                }}
            >
                <defs>
                    <filter id="glow-complete" x="-50%" y="-50%" width="200%" height="200%">
                        <feGaussianBlur in="SourceGraphic" stdDeviation="2" result="blur" />
                        <feMerge>
                            <feMergeNode in="blur" />
                            <feMergeNode in="SourceGraphic" />
                        </feMerge>
                    </filter>
                </defs>

                <motion.circle
                    cx="154"
                    cy="154"
                    r={circleRadius}
                    variants={backgroundCircleVariants1}
                    initial="initial"
                    animate={variant}
                    style={{
                        fill: 'none',
                        strokeWidth: '1px',
                    }}
                    filter={complete ? 'url(#glow-complete)' : ''}
                />
                <motion.circle
                    cx="154"
                    cy="154"
                    r={circleRadius}
                    variants={mainCircleVariants}
                    initial="initial"
                    animate={variant}
                    style={{
                        strokeWidth: '4px',
                        fill: 'none',
                    }}
                    filter={progress === 100 ? 'url(#glow-complete)' : ''}
                />

                <motion.path
                    d="M146.699,121.137l-0,-9.078l-29.71,-0l8.537,9.078l21.173,-0Z"
                    variants={getProgressPathAnimation(72, progress)}
                    initial="initial"
                    animate={variant}
                    style={{
                        stroke: '#e3e3e3',
                        strokeWidth: '2px',
                        fill: '#e3e3e3',
                    }}
                    filter={progress === 100 ? 'url(#glow-complete)' : ''}
                />
                <motion.path
                    d="M156.602,244.103l8.253,0l-0,-75.1l-8.164,8.402l-0.089,66.698Z"
                    variants={getProgressPathAnimation(161, progress)}
                    initial="initial"
                    animate={variant}
                    style={{
                        stroke: '#e3e3e3',
                        strokeWidth: '2px',
                        fill: '#e3e3e3',
                    }}
                    filter={complete ? 'url(#glow-complete)' : ''}
                />
                <motion.path
                    d="M178.059,230.899l8.253,-8.485l0,-64.14l-8.253,0l0,72.625Z"
                    variants={getProgressPathAnimation(156, progress)}
                    initial="initial"
                    animate={variant}
                    style={{
                        stroke: '#e3e3e3',
                        strokeWidth: '2px',
                        fill: '#e3e3e3',
                    }}
                    filter={complete ? 'url(#glow-complete)' : ''}
                />
                <motion.path
                    d="M199.517,209.441l8.253,-8.515l-0,-26.74l-8.253,-8.484l-0,43.739Z"
                    variants={getProgressPathAnimation(94, progress)}
                    initial="initial"
                    animate={variant}
                    style={{
                        stroke: '#e3e3e3',
                        strokeWidth: '2px',
                        fill: '#e3e3e3',
                    }}
                    filter={progress === 100 ? 'url(#glow-complete)' : ''}
                />
                <motion.path
                    d="M207.77,162.022l-12.303,-24.38l-28.452,-0l-31.87,32.016l-0,52.746l8.456,8.495l-0,-57.714l26.925,-27.048l21.43,-0l15.814,15.885Z"
                    variants={getProgressPathAnimation(305, progress)}
                    initial="initial"
                    animate={variant}
                    style={{
                        stroke: '#e3e3e3',
                        strokeWidth: '2px',
                        fill: '#e3e3e3',
                    }}
                    filter={complete ? 'url(#glow-complete)' : ''}
                />
                <motion.path
                    d="M134.728,137.642l-21.04,21.154l-0,42.215l8.52,8.43l-0,-47.242l24.491,-24.557l-11.971,-0Z"
                    variants={getProgressPathAnimation(175, progress)}
                    initial="initial"
                    animate={variant}
                    style={{
                        stroke: '#e3e3e3',
                        strokeWidth: '2px',
                        fill: '#e3e3e3',
                    }}
                    filter={progress === 100 ? 'url(#glow-complete)' : ''}
                />
                <motion.path
                    d="M100.483,187.984l0,-15.68l-8.253,7.855l8.253,7.825Z"
                    variants={getProgressPathAnimation(38, progress)}
                    initial="initial"
                    animate={variant}
                    style={{
                        stroke: '#e3e3e3',
                        strokeWidth: '2px',
                        fill: '#e3e3e3',
                    }}
                    filter={complete ? 'url(#glow-complete)' : ''}
                />
                <motion.path
                    d="M115.497,93.61l39.893,9.756l18.837,19.236l-4.316,4.312l12.116,-0.123l4.285,-4.189l-18.189,-18.594l14.49,-14.373l-22.321,-22.142l-59.809,0l17.295,17.157l-11.314,-0l9.033,8.96Z"
                    variants={getProgressPathAnimation(278, progress)}
                    initial="initial"
                    animate={variant}
                    style={{
                        stroke: '#e3e3e3',
                        strokeWidth: '2px',
                        fill: '#e3e3e3',
                    }}
                    filter={complete ? 'url(#glow-complete)' : ''}
                />

                {/* Background Paths */}
                <path
                    d="M146.699,121.137l-0,-9.078l-29.71,-0l8.537,9.078l21.173,-0Z"
                    style={backgroundPath}
                />
                <path
                    d="M156.602,244.103l8.253,0l-0,-75.1l-8.164,8.402l-0.089,66.698Z"
                    style={backgroundPath}
                />
                <path
                    d="M178.059,230.899l8.253,-8.485l0,-64.14l-8.253,0l0,72.625Z"
                    style={backgroundPath}
                />
                <path
                    d="M199.517,209.441l8.253,-8.515l-0,-26.74l-8.253,-8.484l-0,43.739Z"
                    style={backgroundPath}
                />
                <path
                    d="M207.77,162.022l-12.303,-24.38l-28.452,-0l-31.87,32.016l-0,52.746l8.456,8.495l-0,-57.714l26.925,-27.048l21.43,-0l15.814,15.885Z"
                    style={backgroundPath}
                />
                <path
                    d="M134.728,137.642l-21.04,21.154l-0,42.215l8.52,8.43l-0,-47.242l24.491,-24.557l-11.971,-0Z"
                    style={backgroundPath}
                />
                <path
                    d="M100.483,187.984l0,-15.68l-8.253,7.855l8.253,7.825Z"
                    style={backgroundPath}
                />
                <path
                    d="M115.497,93.61l39.893,9.756l18.837,19.236l-4.316,4.312l12.116,-0.123l4.285,-4.189l-18.189,-18.594l14.49,-14.373l-22.321,-22.142l-59.809,0l17.295,17.157l-11.314,-0l9.033,8.96Z"
                    style={backgroundPath}
                />
            </svg>
        </div>
    );
}
