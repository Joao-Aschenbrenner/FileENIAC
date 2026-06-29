// SPDX-License-Identifier: MIT
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Modal } from '../Modal';
import { describe, it, expect, vi } from 'vitest';

describe('Modal', () => {
  it('does not render when closed', () => {
    render(<Modal open={false} onClose={() => {}} title="Modal">content</Modal>);
    expect(screen.queryByText('content')).not.toBeInTheDocument();
  });

  it('renders children when open', () => {
    render(<Modal open={true} onClose={() => {}} title="Modal">content</Modal>);
    expect(screen.getByText('content')).toBeInTheDocument();
  });

  it('renders title when open', () => {
    render(<Modal open={true} onClose={() => {}} title="Meu Modal">content</Modal>);
    expect(screen.getByText('Meu Modal')).toBeInTheDocument();
  });

  it('calls onClose on escape key', async () => {
    const onClose = vi.fn();
    render(<Modal open={true} onClose={onClose} title="Modal">content</Modal>);
    await userEvent.keyboard('{Escape}');
    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
