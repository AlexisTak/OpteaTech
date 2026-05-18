export type ApiEnvelope<T> = {
  data: T;
  timestamp?: string;
};

export type RequestStatus =
  | 'nouveau'
  | 'en_etude'
  | 'devis_envoye'
  | 'accepte'
  | 'en_cours'
  | 'en_revision'
  | 'livre'
  | 'termine'
  | 'annule';
