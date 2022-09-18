import { FunctionComponent } from 'react';

type TButtonProps = HTMLButtonElement & {label: string};

export const PrimaryButton: FunctionComponent<Partial<TButtonProps>> = ({
    children, label,
    className
}) => {
    return (
        <button className={`bg-black border border-blue-500 hover:opacity-75 transition duration-500 text-white px-3 py-1 text-sm rounded-md ${className}`}>
            {label ?? children as any}
        </button>
    );
}