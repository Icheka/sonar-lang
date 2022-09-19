import { Context, useContext } from 'react';

export const ContextHookFactory = <T>(context: Context<T>, check?: string) => {
    const c = useContext(context);
    if (!c || (check && !(c as Record<string,any>).hasOwnProperty(check))) throw new Error(`${context.displayName}() can only be used within ${context.displayName}Provider`);
    return c;
}

export * from './editor';