import React from "react";

export default function Profile({ name, headline, summary, skills, projects }) {
  return (
    <div className="profile">
      <h1>{name}</h1>
      <h2>{headline}</h2>
      <p>{summary}</p>
      <h3>Skills</h3>
      <ul>{skills.map((s) => <li key={s}>{s}</li>)}</ul>
      <h3>Projects</h3>
      <ul>
        {projects.map((proj) => (
          <li key={proj.title}>
            <a href={proj.link}>{proj.title}</a>: {proj.description}
          </li>
        ))}
      </ul>
    </div>
  );
}
