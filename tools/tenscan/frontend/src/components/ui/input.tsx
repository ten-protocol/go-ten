import styles from "./styles.module.scss"

import { cn } from "@/src/lib/utils";
import {forwardRef} from "react";

export interface InputProps
    extends React.InputHTMLAttributes<HTMLInputElement> {}

const Input = forwardRef<HTMLInputElement, InputProps>(
    ({ className, type, ...props }, ref) => {
        return (
            <div className={styles.input_container}>
                <div className={styles.focus_effect}/>
                <input
                    type={type}
                    className={cn(
                        "flex h-10 w-full border border-input bg-background px-3 py-2 text-sm file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground outline-none disabled:cursor-not-allowed disabled:opacity-50",
                        className,
                    )}
                    ref={ref}
                    {...props}
                />
            </div>
        );
    },
);
Input.displayName = "Input";

export { Input };