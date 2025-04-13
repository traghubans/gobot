import React, { useState, KeyboardEvent } from 'react';
import styled from 'styled-components';

const Container = styled.div`
  padding: 15px;
  background-color: #f0f0f0;
  border-top: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
`;

const Input = styled.textarea`
  flex: 1;
  padding: 10px 15px;
  border-radius: 20px;
  border: none;
  background-color: white;
  font-size: 16px;
  line-height: 1.4;
  resize: none;
  max-height: 100px;
  min-height: 40px;
  font-family: inherit;

  &:focus {
    outline: none;
  }
`;

const SendButton = styled.button<{ hasContent: boolean }>`
  background: none;
  border: none;
  padding: 10px 15px;
  margin-left: 10px;
  color: ${props => props.hasContent ? '#007AFF' : '#999'};
  font-size: 16px;
  font-weight: 600;
  cursor: ${props => props.hasContent ? 'pointer' : 'default'};
  transition: color 0.2s;

  &:focus {
    outline: none;
  }
`;

interface InputBarProps {
  onSendMessage: (message: string) => void;
}

const InputBar: React.FC<InputBarProps> = ({ onSendMessage }) => {
  const [message, setMessage] = useState('');

  const handleSend = () => {
    if (message.trim()) {
      onSendMessage(message);
      setMessage('');
    }
  };

  const handleKeyPress = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(e.target.value);
    e.target.style.height = 'auto';
    e.target.style.height = `${Math.min(e.target.scrollHeight, 100)}px`;
  };

  return (
    <Container>
      <Input
        value={message}
        onChange={handleChange}
        onKeyPress={handleKeyPress}
        placeholder="Type a message..."
      />
      <SendButton
        onClick={handleSend}
        hasContent={message.trim().length > 0}
      >
        Send
      </SendButton>
    </Container>
  );
};

export default InputBar; 