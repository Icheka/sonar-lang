import { FunctionComponent, useEffect, useMemo } from 'react';

import { EditorContextProvider } from '../../context';
import { useEditor } from '../../context/editor';
import { NavigationBar } from '../layout';
import { Editor } from './Editor';

export const Scaffold: FunctionComponent = () => {
    return (
        <div className="h-full">
            <EditorContextProvider>
                <Child />
            </EditorContextProvider>
        </div>  
    );
}

const Child: FunctionComponent = () => {
    // vars
    const {stderr, stdout} = useEditor();

    // utils
    const hydrateOutput = useMemo(() => {
        const out = stderr || stdout;
        return out.replaceAll('\n', '<br />');
    }, [stderr, stdout]);

    // hooks
    const css = useMemo(() => {
        if (stderr.trim().length > 0) return `text-red-500`;
        return `text-white`;
    }, [stderr, stdout]);

    return (
        <div className='flex flex-col w-full h-full'>
            <NavigationBar />
            <div className='overflow-y-auto flex-1'>
                <Editor />
            </div>
            <div dangerouslySetInnerHTML={{__html: hydrateOutput}} className={`h-[18%] bg-gray-600 overflow-y-auto px-3 py-2 text-sm font-mono ${css}`} />
        </div>
    );
}