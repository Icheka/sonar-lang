import 'prismjs/themes/prism.css';

import { highlight, languages } from 'prismjs';
import { useLayoutEffect, useMemo, useState } from 'react';
import { FunctionComponent } from 'react';
import EditorCore from 'react-simple-code-editor';

export const Editor: FunctionComponent = () => {
    // vars
    const textAreaClassName = 'editor-core-textarea';

    // state
    const [code, setCode] = useState('');
    const [domPainted, setDomPainted] = useState(false);

    // hooks
    useLayoutEffect(() => {
        setDomPainted(true);
    }, []);
    const textArea = useMemo(() => {
        return document.querySelector(`.${textAreaClassName}`) as HTMLTextAreaElement;
    }, [textAreaClassName, domPainted]);
    const lineNumbers = useMemo(() => {
        return code
            .split('\n')
            .map((_, i) => String(i + 1))
            .join('\n');
    }, [code]);

    // utils
    const selectTextarea = () => {
        if (!textArea) return;
        textArea.focus();
    }

    return (
        <>
            <div onClick={selectTextarea} className="bg-red-200 w-full h-full overflow-y-auto">
                <section className='overflow-y-auto flex'>
                    <pre className='bg-blue-200 text-xs leading-[18px] pt-[11px] pb-4 px-2 opacity-50'>{lineNumbers}</pre>
                    <div className='w-full'>
                        <EditorCore
                            value={code}
                            onValueChange={setCode}
                            highlight={code => highlight(code, languages.js, 'js')}
                            padding={10}
                            style={{
                                fontFamily: '"Fira code", "Fira Mono", monospace',
                                fontSize: 12,
                            }}
                            tabSize={4}
                            insertSpaces
                            preClassName='editor-core-pre'
                            textareaClassName={textAreaClassName}
                            className="editor-core"
                        />
                    </div>
                </section>
            </div>  
        </>
    );
}