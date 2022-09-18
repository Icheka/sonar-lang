import { FunctionComponent, useEffect } from 'react';

import { EditorContextProvider } from '../../context';
import { NavigationBar } from '../layout';
import { Editor } from './Editor';

export const Scaffold: FunctionComponent = () => {
    return (
        <div className="bg-blue-900 h-full">
            <EditorContextProvider>
                <NavigationBar />
                <div className='rounded-md overflow-y-auto border-2 border-indigo-400 h-full max-h-[94vh]'>
                    <Editor />
                </div>
            </EditorContextProvider>
        </div>  
    );
}