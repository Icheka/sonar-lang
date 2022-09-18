import { createContext, Dispatch, FunctionComponent, SetStateAction, useContext, useMemo, useState } from 'react';

import { ContextHookFactory } from '.';

interface IEditorContext {
    code: string;
    setCode: Dispatch<SetStateAction<string>>;
    stderr: string;
    stdout: string;
}

const EditorContext = createContext<IEditorContext>({} as IEditorContext);
EditorContext.displayName = 'EditorContext';

export const EditorContextProvider: FunctionComponent<{children: any}> = ({children}) => {
    // state
    const [code, setCode] = useState('');
    const [stderr, setStderr] = useState('');
    const [stdout, setStdout] = useState('');

    // hooks
    const value = useMemo(() => ({
        code,
        setCode,
        stderr,
        stdout
    }), [code, stderr, stdout]);

    return (
        <EditorContext.Provider value={value}>
            {children}
        </EditorContext.Provider>
    );
}

export const useEditor = () => ContextHookFactory(EditorContext, 'code');