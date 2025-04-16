import React from 'react';

// Стили для анимации раскрытия
const slideDownAnimation = `
  @keyframes slideDown {
    from {
      max-height: 0;
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      max-height: 1000px;
      opacity: 1;
      transform: translateY(0);
    }
  }
  
  .animate-slideDown {
    animation: slideDown 0.3s ease-out forwards;
  }
`;

// CSS-анимация пульсации для кнопок раскрытия
const pulseAnimation = `
  @keyframes pulse {
    0% {
      box-shadow: 0 0 0 0 rgba(59, 130, 246, 0.4);
    }
    70% {
      box-shadow: 0 0 0 6px rgba(59, 130, 246, 0);
    }
    100% {
      box-shadow: 0 0 0 0 rgba(59, 130, 246, 0);
    }
  }
`;

// Компонент, который добавляет все CSS-стили
const TaskStyles: React.FC = () => {
  const combinedStyles = `
    ${slideDownAnimation}
    ${pulseAnimation}
  `;
  
  return <style>{combinedStyles}</style>;
};

export default TaskStyles; 