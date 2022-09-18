import { Editor } from './components/domains';
import { NavigationBar } from './components/layout';

function App() {
    return (
        <div className="h-[94vh] bg-blue-900">
            <NavigationBar />
            <div className='rounded-md overflow-y-auto border-2 border-indigo-400 h-full max-h-[94vh]'>
                <Editor />
            </div>
        </div>  
    );
}

export default App;
