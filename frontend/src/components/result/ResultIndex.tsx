import React, { useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axiosInstance           from '@/axiosConfig';
import { ResultSummary }       from '@/types/result';
import '@/styles/components/result/_common.css';
import '@/styles/components/result/ResultIndex.css';

/* 品詞 ID → 名称の簡易マップ（必要に応じて拡張） */
const POS_MAP: Record<number, string> = {
  1: '名詞',
  2: '代名',
  3: '動詞',
  4: '形容',
  5: '副詞',
  // …
}

/* ページサイズ選択肢 */
const PAGE_SIZES = [10, 20, 30] as const

const ResultIndex: React.FC = () => {
  const nav = useNavigate()

  /* ---------- state ---------- */
  const [list, setList] = useState<ResultSummary[]>([])
  const [pageSize, setPageSize] =
    useState<(typeof PAGE_SIZES)[number]>(10)
  const [page, setPage] = useState(0) // 0-based
  const [loading, setLoading] = useState(true)
  const [errMsg, setErrMsg] = useState('')

  /* ---------- 初回ロード ---------- */
  useEffect(() => {
    const fetchAll = async () => {
      try {
        const res = await axiosInstance.get<ResultSummary[]>('/results')
        setList(
          res.data.sort(
            (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
          ),
        )
      } catch (e) {
        console.error(e)
        setErrMsg('成績の取得に失敗しました')
      } finally {
        setLoading(false)
      }
    }
    fetchAll()
  }, [])

  /* ---------- ページ計算 ---------- */
  const rows = useMemo(() => {
    const start = page * pageSize
    return list.slice(start, start + pageSize)
  }, [list, page, pageSize])
 /* ---------- UI ---------- */
 if (loading) return <p className="p-4">読み込み中…</p>;
 if (errMsg)  return <p className="p-4 text-red-600">{errMsg}</p>;

 return (
   <div className="rs-wrapper">
     <h1 className="ri-title">成績一覧</h1>

     {/* --- 一覧テーブル --- */}
     <table className="rs-table ri-table">
       <thead>
         <tr>
           <th></th><th>日付</th><th>登録単語</th><th>慣用句</th>
           <th>特殊</th><th>品詞</th><th>問題</th><th>正解</th><th>正解率</th>
         </tr>
       </thead>
       <tbody>
         {rows.map(r=>(
           <tr key={r.quizNumber} onClick={()=>nav(`/results/${r.quizNumber}`)}
               className="cursor-pointer">
             {/* 以下セルは以前と同じ */}
             <td>{r.quizNumber}</td>
             <td className="whitespace-nowrap">{new Date(r.createdAt).toLocaleString()}</td>
             <td>{['全','登録','未登録'][r.isRegisteredWords]}</td>
             <td>{['全','のみ','除外'][r.isIdioms]}</td>
             <td>{['全','のみ','除外'][r.isSpecialCharacters]}</td>
             <td>{r.choicesPosIds.map(id=>POS_MAP[id]??id).join(', ')}</td>
             <td>{r.totalQuestionsCount}</td>
             <td>{r.correctCount}</td>
             <td>{r.resultCorrectRate.toFixed(1)}%</td>
           </tr>
         ))}
       </tbody>
     </table>

     {/* --- ページネーション --- */}
     <div className="flex items-center justify-between">
       {/* サイズボタン */}
       <div className="space-x-2">
         {PAGE_SIZES.map(s=>(
           <button key={s}
             className={`rs-page-btn ${s===pageSize?'active':''}`}
             onClick={()=>{setPageSize(s);setPage(0);}}>
             {s}
           </button>
         ))}
       </div>

       {/* Prev / Next */}
       <div className="space-x-2">
         <button onClick={()=>setPage(p=>p-1)} disabled={page===0}
                 className="rs-page-btn">Prev</button>
         <button onClick={()=>setPage(p=>p+1)}
                 disabled={(page+1)*pageSize>=list.length}
                 className="rs-page-btn">Next</button>
       </div>
     </div>
   </div>
 );
};

export default ResultIndex;