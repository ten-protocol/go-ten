'use client';

import { useScramble } from 'use-scramble';

type Props = {
    text: string;
    hover?: boolean;
    active?: boolean;
    onView?: boolean;
};

export default function EncryptedTextAnimation({
    text,
    hover = true,
}: Props) {
    const { ref, replay } = useScramble({
        text,
        overflow: true,
        overdrive: false,
        playOnMount: false,
        ignore: [' '],
        range: [35, 38, 48, 57, 97, 122],
        speed: 0.6,
        tick: 1,
        step: 1,
        scramble: 10,
        seed: 5,
        chance: 0.9,
    });

    const handleInteraction = () => {
        if (hover) {
            replay();
        }
    };

    return (
        <p ref={ref} onMouseOver={handleInteraction} onFocus={handleInteraction}>
            {text}
        </p>
    );
}
