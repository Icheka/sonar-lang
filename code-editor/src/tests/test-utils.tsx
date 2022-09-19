import { render, RenderOptions } from '@testing-library/react';
import { FunctionComponent, ReactElement, ReactNode } from 'react';

import { EditorContextProvider } from '../context';

const AllProviders: FunctionComponent<{children: ReactNode}> = ({children}) => (
    <EditorContextProvider>
        {children}
    </EditorContextProvider>
);

const customRender = (
    ui: ReactElement,
    options?: Omit<RenderOptions, 'wrapper'>
) => render(ui, {wrapper: AllProviders, ...options});

export * from '@testing-library/react';
export {customRender as render};