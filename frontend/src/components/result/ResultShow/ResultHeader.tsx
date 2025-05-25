import React from 'react';

type Props = { correct: number; total: number; rate: number };

const ResultHeader: React.FC<Props> = ({ correct, total, rate }) => (
  <h1 className="rs-summary">
    {correct}/{total} 正解
    <span>({rate.toFixed(1)}%)</span>
  </h1>
);

export default ResultHeader;
