import { FunctionComponent, MouseEventHandler } from 'react';

type TButtonProps = HTMLButtonElement & {
    label: string,
    onClick: MouseEventHandler<HTMLButtonElement>;
};

export const PrimaryButton: FunctionComponent<Partial<TButtonProps>> = ({
    children, label,
    className, onClick
}) => {
    return (
        <button onClick={(e) => onClick && onClick(e)} className={`bg-black border border-blue-500 hover:opacity-75 transition duration-500 text-white px-3 py-1 text-sm rounded-md ${className}`}>
            {label ?? children as any}
        </button>
    );
}