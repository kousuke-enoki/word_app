import clsx from 'clsx';

import { ResultQuestion } from '@/types/quiz';

type Props = {
  rows:            ResultQuestion[];
  onToggleRegister:(row: ResultQuestion)=>void;
  onChangeAtten:   (row: ResultQuestion, lv: number)=>void;
};

const ResultTable = ({ rows, onToggleRegister, onChangeAtten }: Props) => (
  <table className="rs-table">
    <thead>
      <tr>
        <th></th><th>単語</th><th>正解</th><th>選択</th>
        <th>登録</th><th>回答回数</th><th>正解回数</th><th>注意Lv.</th>
      </tr>
    </thead>
    <tbody>
      {rows.map(row=>{
        const correct = row.choicesJpms.find(c=>c.japaneseMeanID===row.correctJpmID)?.name ?? '-';
        const select  = row.choicesJpms.find(c=>c.japaneseMeanID===row.answerJpmID )?.name ?? '-';
        return (
          <tr key={row.questionNumber}>
            <td>{row.questionNumber}</td>
            <td>{row.wordName}</td>
            <td>{correct}</td>
            <td className={clsx(
                'border p-1',
                row.isCorrect ? 'correct-cell' : 'wrong-cell'
              )}
            >{select}</td>
            <td>
              <button onClick={()=>onToggleRegister(row)}
                      className={row.registeredWord.isRegistered
                        ? 'bg-red-600 hover:bg-red-700 text-white px-2 py-1 rounded'
                        : 'bg-green-600 hover:bg-green-700 text-white px-2 py-1 rounded'}>
                {row.registeredWord.isRegistered?'解除':'登録'}
              </button>
            </td>
            <td>{row.registeredWord.quizCount}</td>
            <td>{row.registeredWord.correctCount}</td>
            <td>
              <select defaultValue={row.registeredWord.attentionLevel}
                      onChange={e=>onChangeAtten(row, +e.target.value)}>
                {[1,2,3,4,5].map(lv=><option key={lv}>{lv}</option>)}
              </select>
            </td>
          </tr>
        );
      })}
    </tbody>
  </table>
);

export default ResultTable;
