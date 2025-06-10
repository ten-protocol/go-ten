import Link from 'next/link';
import { FaDiscord, FaGithub, FaTelegram } from 'react-icons/fa';
import { RiTwitterXFill } from 'react-icons/ri';

export default function Social() {
    return (
        <div className="flex gap-6">
            <Link href="https://twitter.com/tenprotocol" target="_blank">
                <RiTwitterXFill size={24} />
            </Link>
            <Link href="https://discord.gg/tenprotocol" target="_blank">
                <FaDiscord size={24} />
            </Link>
            <Link href="https://t.me/tenprotocol" target="_blank">
                <FaTelegram size={24} />
            </Link>
            <Link href="https://github.com/ten-protocol" target="_blank">
                <FaGithub size={24} />
            </Link>
        </div>
    );
}
