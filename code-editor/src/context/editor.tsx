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
    ping: 'ping',
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
    const connOpen = () => socket.readyState === socket.OPEN;
    const emit = (payload: SocketReqRes) => {
        if (!connOpen()) {
            socket = wsConnect();
        }
        const interval = setInterval(() => {
            if (connOpen()) {
                clearInterval(interval);
                socket.send(JSON.stringify(payload));
            }
        }, 2000);
    }
    const evaluateCode = () => {
        emit({
            type: 'evaluate',
            data: code
        });
    }
    const clearStd = () => {
        setStderr("");
        setStdout("");
    }
    const handleStdErr = (data: string) => {
        clearStd();
        setStderr(data);
    }
    const handleStdOut = (data: string) => {
        clearStd();
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
        socket = wsConnect();

        socket.onopen = () => {
            console.log('Connected.');
            setConnected(true);

            setInterval(() => {
                socket.send(JSON.stringify({type: "ping"}));
            }, 50000);
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
    }, [connOpen()]);

    return (
        <EditorContext.Provider value={value}>
            {children}
        </EditorContext.Provider>
    );
}

export const useEditor = () => ContextHookFactory(EditorContext, 'code');