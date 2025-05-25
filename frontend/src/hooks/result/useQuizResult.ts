import { useEffect, useReducer } from 'react';
import axiosInstance from '@/axiosConfig';
import { ResultRes } from '@/types/quiz';

type State =
  | { loading: true;  error: null; result: null }
  | { loading: false; error: null; result: ResultRes }
  | { loading: false; error: Error; result: null };

type Action =
  | { type: 'START' }
  | { type: 'SUCCESS'; payload: ResultRes }
  | { type: 'ERROR';  payload: Error };

const reducer = (_: State, a: Action): State => {
  switch (a.type) {
    case 'START':   return { loading: true,  error: null,  result: null };
    case 'SUCCESS': return { loading: false, error: null,  result: a.payload };
    case 'ERROR':   return { loading: false, error: a.payload, result: null };
  }
};

export const useQuizResult = (quizNo?: string) => {
  const [state, dispatch] = useReducer(reducer, {
    loading: true,
    error:   null,
    result:  null,
  } as State);

  useEffect(() => {
    if (!quizNo) return;

    const abort = new AbortController();
    (async () => {
      try {
        dispatch({ type: 'START' });
        const { data } = await axiosInstance.get<ResultRes>(
          `/results/${quizNo}`,
          { signal: abort.signal },
        );
        dispatch({ type: 'SUCCESS', payload: data });
      } catch (e: any) {
        if (!abort.signal.aborted) dispatch({ type: 'ERROR', payload: e });
      }
    })();
    return () => abort.abort();
  }, [quizNo]);

  return state;
};
