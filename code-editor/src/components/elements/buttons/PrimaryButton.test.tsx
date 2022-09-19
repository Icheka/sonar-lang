import { fireEvent, render, screen } from '../../../tests';
import { PrimaryButton } from './PrimaryButton';

test('renders an interactive button element', async () => {
    let counter = 0;
    render(<PrimaryButton label="Click me" onClick={() => counter++} />);

    const element = screen.getByText('Click me');
    expect(element.tagName).toBe('BUTTON');

    fireEvent.click(element);
    expect(counter).toBe(1);
});