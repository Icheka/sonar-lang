import { FunctionComponent } from 'react';

import { PrimaryButton } from '../elements';

export const NavigationBar: FunctionComponent = () => {
    return (
        <nav className='flex justify-between items-center bg-gray-700 py-2 px-3'>
            <span>
            [LOGO]
            </span>
            
            <div className="flex items-center space-x-4">
                <PrimaryButton label="Run" />
            </div>
        </nav>
    )
}