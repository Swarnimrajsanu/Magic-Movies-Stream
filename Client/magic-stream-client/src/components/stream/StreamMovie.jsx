import { useParams } from "react-router-dom";
import "./StreamMovie.css";

const StreamMovie = () => {
  const { yt_id } = useParams();

  if (!yt_id) return <h2>Video not found</h2>;

  return (
    <div className="stream-container">
      <iframe
        className="stream-player"
        src={`https://www.youtube-nocookie.com/embed/${yt_id}?autoplay=1`}
        title="Movie Stream"
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; fullscreen"
        allowFullScreen
      ></iframe>
    </div>
  );
};

export default StreamMovie;
