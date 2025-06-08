const PAGE_SIZES = [10, 30, 50, 100] as const;
type PageSize = typeof PAGE_SIZES[number];

type Props = {
  sizes?:  readonly PageSize[];
  size:    number;
  page:    number;
  total:   number;
  onSize:  (s:PageSize)=>void;
  onPrev:  ()=>void;
  onNext:  ()=>void;
};

const Pagination = ({ sizes = PAGE_SIZES, size, page, total, onSize, onPrev, onNext }:Props) => (
  <div className="flex items-center justify-between">
    <div className="rs-page-size space-x-2">
      {sizes.map(s=>(
        <button key={s}
                className={`rs-page-btn ${s===size?'active':''}`}
                onClick={()=>onSize(s)}>{s}</button>
      ))}
    </div>
    <div className="rs-pager space-x-2">
      <button disabled={page===0}               onClick={onPrev}>Prev</button>
      <button disabled={(page+1)*size>=total}   onClick={onNext}>Next</button>
    </div>
  </div>
);

export default Pagination;
