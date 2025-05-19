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

// Стили для блока расходов
const expensesStyles = `
  .expenses-list {
    margin-top: 0.5rem;
    transition: all 0.2s ease;
  }
  
  .expense-item {
    display: flex;
    align-items: center;
    padding: 0.25rem 0.5rem;
    border-radius: 0.25rem;
    transition: background-color 0.2s ease;
  }
  
  .expense-item:hover {
    background-color: rgba(34, 197, 94, 0.05);
  }
  
  .expenses-summary {
    font-size: 0.75rem;
    color: #065f46;
    font-weight: 500;
    padding: 0.25rem 0.5rem;
    background-color: rgba(34, 197, 94, 0.1);
    border-radius: 0.25rem;
    display: inline-flex;
    align-items: center;
    margin-top: 0.25rem;
  }
  
  /* Стили для отображения иерархии расходов */
  .child-expenses {
    padding-left: 0.5rem;
    margin-left: 0.5rem;
    border-left: 2px solid rgba(34, 197, 94, 0.2);
  }
  
  .child-expense-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 0.5rem;
    padding: 0.25rem;
    border-radius: 0.25rem;
    background-color: rgba(34, 197, 94, 0.05);
  }
  
  /* Анимация для расходов */
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  
  .expenses-animation {
    animation: fadeIn 0.3s ease-in;
  }
`;

// Компонент, который добавляет все CSS-стили
const TaskStyles: React.FC = () => {
  const combinedStyles = `
    ${slideDownAnimation}
    ${pulseAnimation}
    ${expensesStyles}
  `;
  
  return <style>{combinedStyles}</style>;
};

export default TaskStyles; 