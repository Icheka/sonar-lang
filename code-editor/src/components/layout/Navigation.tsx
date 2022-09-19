import { FunctionComponent } from 'react';

import { useEditor } from '../../context/editor';
import { PrimaryButton } from '../elements';

export const NavigationBar: FunctionComponent = () => {
    // vars
    const {evaluate} = useEditor();

    return (
        <nav className='flex justify-between items-center bg-gray-700 py-2 px-3'>
            <span>
            [LOGO]
            </span>
            
            <div className="flex items-center space-x-4">
                <PrimaryButton onClick={evaluate} label="Run" />
            </div>
        </nav>
    )
}