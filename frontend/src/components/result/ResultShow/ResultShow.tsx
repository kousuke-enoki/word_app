import React, { useMemo, useState, useCallback } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuizResult }  from '@/hooks/result/useQuizResult';
import ResultHeader       from '@/components/result/ResultShow/ResultHeader';
import ResultSettingCard  from '@/components/result/ResultShow/ResultSettingCard';
import ResultTable        from '@/components/result/ResultShow/ResultTable';
import Pagination         from '@/components/common/Pagination';
import { ResultQuestion } from '@/types/quiz';
import { registerWord }   from '@/service/word/RegisterWord';
import '@/styles/components/result/_common.css';
import '@/styles/components/result/ResultShow.css';


const pageSizes = [10,30,50,100] as const;
type PageSize = typeof pageSizes[number];

export default function ResultShow() {
  const { quizNo } = useParams<{quizNo?:string}>();
  const nav  = useNavigate();
  const { loading,error,result } = useQuizResult(quizNo);

  /* --- ページング状態 ------------------------------------ */
  const [size,setSize] = useState<PageSize>(10);
  const [page,setPage] = useState(0);

  /** 表示対象行（memo） */
  const rows = useMemo(()=>{
    if (!result) return [];
    const start = page*size;
    return result.resultQuestions.slice(start,start+size);
  },[result,page,size]);

  /* --- 行操作コールバック -------------------------------- */
  const toggleRegister = useCallback( async (row:ResultQuestion)=>{
    try{
      const u = await registerWord(row.wordID,!row.registeredWord.isRegistered);
      row.registeredWord.isRegistered = u.isRegistered;
      row.registeredWord.quizCount    = u.quizCount    ?? row.registeredWord.quizCount;
      row.registeredWord.correctCount = u.correctCount ?? row.registeredWord.correctCount;
      // → Immer や Redux Toolkit を使う場合はここで dispatch
    }catch(e){ console.error(e); }
  },[]);

  const changeAttention = useCallback( async (row:ResultQuestion)=>{
    try{ await registerWord(row.wordID,true); }catch(e){ console.error(e); }
  },[]);

  /* --- サイズ候補 (問題数以下のみ) ------------------------ */
  const sizeCandidates: PageSize[] = useMemo(()=>{
    if (!result) return [];
    return pageSizes.filter((s):s is PageSize=>s<=result.totalQuestionsCount);
  },[result]);

  /* --- UI ----------------------------------------------- */
  if (loading) return <p className="text-center py-10">Loading...</p>;
  if (error   ) return <p className="text-center py-10 text-red-500">通信に失敗しました</p>;
  if (!result ) return null;

  return (
    <div className="rs-wrapper">
      <ResultHeader
        correct={result.correctCount}
        total={result.totalQuestionsCount}
        rate={result.resultCorrectRate}
      />

      <ResultSettingCard setting={result.resultSetting} />

      <ResultTable
        rows={rows}
        onToggleRegister={toggleRegister}
        onChangeAtten={changeAttention}
      />

      <Pagination
        sizes={sizeCandidates}
        size={size}
        page={page}
        total={result.resultQuestions.length}
        onSize={(s:PageSize)=>{ setSize(s); setPage(0); }}
        onPrev={()=>setPage(p=>p-1)}
        onNext={()=>setPage(p=>p+1)}
      />

      {/* 下部ナビだけはこのコンポーネント内に残す */}
      <div className="rs-nav mt-6 flex justify-center gap-4">
        <button onClick={()=>nav('/quizs')}   className="primary">クイズメニュー</button>
        <button onClick={()=>nav('/results')} className="primary">成績一覧</button>
        <button onClick={()=>nav('/')}        className="secondary">ホーム</button>
      </div>
    </div>
  );
}
