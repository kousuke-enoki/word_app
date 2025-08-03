import '@/styles/components/result/_common.css';
import '@/styles/components/result/ResultShow.css';

import React, { useCallback,useEffect,useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import Pagination         from '@/components/common/Pagination';
import ResultHeader       from '@/components/result/ResultShow/ResultHeader';
import ResultSettingCard  from '@/components/result/ResultShow/ResultSettingCard';
import ResultTable        from '@/components/result/ResultShow/ResultTable';
import { useQuizResult }  from '@/hooks/result/useQuizResult';
import { registerWord }   from '@/service/word/RegisterWord';
import { ResultQuestion } from '@/types/quiz';


const pageSizes = [10,30,50,100] as const;
type PageSize = typeof pageSizes[number];

export default function ResultShow() {
  const { quizNo } = useParams<{quizNo?:string}>();
  const nav  = useNavigate();
  const { loading,error,result } = useQuizResult(quizNo);

  /* --- ページング状態 ------------------------------------ */
  const [size,setSize] = useState<PageSize>(10);
  const [page,setPage] = useState(0);

  // 表示用のローカル state を持つ
  const [view, setView] = useState<typeof result | null>(null);

  // 初回取得 or quizNo 変化時に view を同期
  useEffect(() => {
    setView(result ?? null);
  }, [result]);

  /** 表示対象行（memo） */
  const rows = useMemo(()=>{
    if (!view) return [];
    const start = page*size;
    return view.resultQuestions.slice(start, start+size);
  },[view, page, size]); // ← result ではなく view を依存に

  /* --- 行操作コールバック -------------------------------- */
  const toggleRegister = useCallback( async (row:ResultQuestion)=>{
    try{
      const u = await registerWord(row.wordID,!row.registeredWord.isRegistered);
      setView(prev => {
        if (!prev) return prev;
        return {
          ...prev,
          resultQuestions: prev.resultQuestions.map(q =>
            q.wordID === row.wordID
              ? {
                  ...q,
                  registeredWord: {
                    ...q.registeredWord,
                    isRegistered: u.isRegistered,
                    quizCount:    u.quizCount    ?? q.registeredWord.quizCount,
                    correctCount: u.correctCount ?? q.registeredWord.correctCount,
                  },
                }
              : q
          ),
        };
      });
      console.log(u)
    }catch(e){ console.error(e); }
  },[]);

  // const changeAttention = useCallback( async (row:ResultQuestion)=>{
  //   try{ await registerWord(row.wordID,true); }catch(e){ console.error(e); }
  // },[]);

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
        // onChangeAtten={changeAttention}
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
