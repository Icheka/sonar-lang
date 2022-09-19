import { createContext, Dispatch, FunctionComponent, SetStateAction, useContext, useEffect, useMemo, useState } from 'react';

import { ContextHookFactory } from '.';
import { keys } from '../config';

interface IEditorContext {
    code: string;
    setCode: Dispatch<SetStateAction<string>>;
    stderr: string;
    stdout: string;
    evaluate: VoidFunction;
}

const EditorContext = createContext<IEditorContext>({} as IEditorContext);
EditorContext.displayName = 'EditorContext';

const wsConnect = () => {
    if (!keys.languageServerURL) throw new Error(`keys.languageServerURL not defined`);
    return new WebSocket(keys.languageServerURL);
}
let socket = wsConnect();
const events = {
    connect: 'connect',
    evaluate: 'evaluate',
    stderr: 'stderr',
    stdout: 'stdout',
}
export type SocketReqRes = {
    type: keyof typeof events;
    data: any;
}

export const EditorContextProvider: FunctionComponent<{children: any}> = ({children}) => {
    // state
    const [code, setCode] = useState('');
    const [stderr, setStderr] = useState('');
    const [stdout, setStdout] = useState('');
    const [connected, setConnected] = useState(false);

    // utils
    const emit = (payload: SocketReqRes) => {
        const connOpen = () => socket && Object.hasOwn(socket, "readyState") && socket.readyState !== socket.OPEN;

        if (!connOpen()) {
            socket = wsConnect();
        }
        const interval = setInterval(() => {
            if (connOpen()) {
                clearInterval(interval);
                socket.send(JSON.stringify(payload));
            }
        }, 1000);
    }
    const evaluateCode = () => {
        console.log("Emitting 'evaluate'");
        emit({
            type: 'evaluate',
            data: code
        });
    }
    const handleStdErr = (data: string) => {
        setStderr(data);
    }
    const handleStdOut = (data: string) => {
        setStdout(data);
    }

    // hooks
    const value: IEditorContext = useMemo(() => ({
        code,
        setCode,
        stderr,
        stdout,
        evaluate: evaluateCode
    }), [code, stderr, stdout]);

    useEffect(() => {
        socket.onopen = () => {
            console.log('Connected.');
            socket.send(JSON.stringify({type: 'connect'}));
            setConnected(true);
        }
        socket.onmessage = e => {
            const data = JSON.parse(e.data) as SocketReqRes;
            switch (data.type) {
                case events.stderr:
                    return handleStdErr(data.data);
                case events.stdout:
                    return handleStdOut(data.data);
            
                default:
                    break;
            }
        }

        return () => {
            socket.close();
        }
    }, []);

    return (
        <EditorContext.Provider value={value}>
            {children}
        </EditorContext.Provider>
    );
}

export const useEditor = () => ContextHookFactory(EditorContext, 'code');