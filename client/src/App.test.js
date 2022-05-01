import { render, screen } from '@testing-library/react';
import App from './App';

test('renders Studio 10B in header', () => {
  render(<App />);
  const brandElement = screen.getByText(/Studio 10B/i);
  expect(brandElement).toBeInTheDocument();
});
